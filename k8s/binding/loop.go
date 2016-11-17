package binding

import (
	"context"
	"errors"

	"github.com/deis/steward-framework"
	"github.com/deis/steward-framework/k8s/data"
	"github.com/deis/steward-framework/k8s/refs"
	"k8s.io/client-go/pkg/api/unversioned"
	apiv1 "k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/pkg/watch"
)

var (
	ErrCancelled   = errors.New("stopped")
	ErrNotABinding = errors.New("not a binding")
	ErrWatchClosed = errors.New("watch closed")
)

// RunLoop starts a blocking control loop that watches and takes action on broker resources.
//
// TODO: remove the namespace param. See https://github.com/deis/steward-framework/issues/30
func RunLoop(
	ctx context.Context,
	namespace string,
	binder framework.Binder,
	secretWriter SecretWriterFunc,
	fn WatchBindingFunc,
	updateFn UpdateBindingFunc,
	getBrokerFn refs.BrokerGetterFunc,
	getSvcClassFn refs.ServiceClassGetterFunc,
	getInstanceFn refs.InstanceGetterFunc,
) error {

	watcher, err := fn(namespace)
	if err != nil {
		return err
	}
	ch := watcher.ResultChan()
	defer watcher.Stop()
	for {
		select {
		case <-ctx.Done():
			return ErrCancelled
		case evt, open := <-ch:
			if !open {
				logger.Errorf("binding loop watch channel was closed")
				return ErrWatchClosed
			}
			logger.Debugf("binding loop received event")
			switch evt.Type {
			case watch.Added:
				if err := handleAddBinding(
					ctx,
					binder,
					updateFn,
					secretWriter,
					getBrokerFn,
					getSvcClassFn,
					getInstanceFn,
					evt,
				); err != nil {
					// TODO: try the handler again. See https://github.com/deis/steward-framework/issues/34
					logger.Errorf("add binding event handler failed (%s)", err)
				}
			}
		}
	}
}

func handleAddBinding(
	ctx context.Context,
	binder framework.Binder,
	updateFn UpdateBindingFunc,
	secretWriter SecretWriterFunc,
	getBrokerFn refs.BrokerGetterFunc,
	getSvcClassFn refs.ServiceClassGetterFunc,
	getInstanceFn refs.InstanceGetterFunc,
	evt watch.Event,
) error {
	binding, ok := evt.Object.(*data.Binding)
	if !ok {
		return ErrNotABinding
	}

	binding.Status.State = data.BindingStatePending
	if _, err := updateFn(binding); err != nil {
		logger.Errorf("error updating binding state to %s", binding.Status.State)
		return err
	}

	broker, serviceClass, instance, err := refs.GetDependenciesForBinding(
		binding,
		getBrokerFn,
		getSvcClassFn,
		getInstanceFn,
	)
	if err != nil {
		logger.Errorf("getting binding %s's dependencies (%s)", binding.Spec.ID, err)
		return err
	}

	bindReq := &framework.BindRequest{
		InstanceID: instance.Spec.ID,
		ServiceID:  serviceClass.ID,
		PlanID:     instance.Spec.PlanID,
		BindingID:  binding.Spec.ID,
		Parameters: binding.Spec.Parameters,
	}
	bindResp, err := binder.Bind(ctx, broker.Spec, bindReq)
	if err != nil {
		logger.Errorf("calling bind operation (%s)", err)
		return err
	}

	secret := &apiv1.Secret{
		TypeMeta: unversioned.TypeMeta{
			Kind: "Secret",
		},
		ObjectMeta: apiv1.ObjectMeta{
			Name:      binding.Spec.SecretName,
			Namespace: binding.Namespace,
		},
		Data: bindResponseCredsToSecretData(bindResp.Creds),
	}

	if _, err := secretWriter(secret); err != nil {
		logger.Errorf("writing secret %s (%s)", binding.Spec.SecretName, err)
		return err
	}

	binding.Status.State = data.BindingStateBound
	if _, err := updateFn(binding); err != nil {
		logger.Errorf("error updating binding state to %s", binding.Status.State)
		return err
	}

	return nil
}

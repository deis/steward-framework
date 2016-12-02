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

// RunLoop starts a blocking control loop that watches and takes action on serviceBroker resources.
//
// TODO: remove the namespace param. See https://github.com/deis/steward-framework/issues/30
func RunLoop(
	ctx context.Context,
	namespace string,
	binder framework.Binder,
	secretWriter SecretWriterFunc,
	fn WatchBindingFunc,
	updateFn UpdateBindingFunc,
	getSvcBrokerFn refs.ServiceBrokerGetterFunc,
	getSvcClassFn refs.ServiceClassGetterFunc,
	getSvcInstanceFn refs.ServiceInstanceGetterFunc,
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
					getSvcBrokerFn,
					getSvcClassFn,
					getSvcInstanceFn,
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
	getSvcBrokerFn refs.ServiceBrokerGetterFunc,
	getSvcClassFn refs.ServiceClassGetterFunc,
	getSvcInstanceFn refs.ServiceInstanceGetterFunc,
	evt watch.Event,
) error {

	binding := new(data.Binding)
	if err := data.TranslateToTPR(evt.Object, binding, data.BindingKind); err != nil {
		return ErrNotABinding
	}

	binding.Status.State = data.BindingStatePending
	binding, err := updateFn(binding)
	if err != nil {
		logger.Errorf("error updating binding state to %s", binding.Status.State)
		return err
	}

	serviceBroker, serviceClass, serviceInstance, err := refs.GetDependenciesForBinding(
		binding,
		getSvcBrokerFn,
		getSvcClassFn,
		getSvcInstanceFn,
	)
	if err != nil {
		logger.Errorf("getting binding %s's dependencies (%s)", binding.Spec.ID, err)
		return err
	}

	bindReq := &framework.BindRequest{
		InstanceID: serviceInstance.Spec.ID,
		ServiceID:  serviceClass.ID,
		PlanID:     serviceInstance.Spec.PlanID,
		BindingID:  binding.Spec.ID,
		Parameters: binding.Spec.Parameters,
	}
	logger.Debugf("issuing binding request %+v", *bindReq)
	bindResp, err := binder.Bind(ctx, serviceBroker.Spec, bindReq)
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

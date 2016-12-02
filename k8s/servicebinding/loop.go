package servicebinding

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
	ErrCancelled          = errors.New("stopped")
	ErrNotAServiceBinding = errors.New("not a service binding")
	ErrWatchClosed        = errors.New("watch closed")
)

// RunLoop starts a blocking control loop that watches and takes action on serviceBroker resources.
//
// TODO: remove the namespace param. See https://github.com/deis/steward-framework/issues/30
func RunLoop(
	ctx context.Context,
	namespace string,
	binder framework.Binder,
	secretWriter SecretWriterFunc,
	fn WatchServiceBindingFunc,
	updateFn UpdateServiceBindingFunc,
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
				logger.Errorf("service binding loop watch channel was closed")
				return ErrWatchClosed
			}
			logger.Debugf("service binding loop received event")
			switch evt.Type {
			case watch.Added:
				if err := handleAddServiceBinding(
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
					logger.Errorf("add service binding event handler failed (%s)", err)
				}
			}
		}
	}
}

func handleAddServiceBinding(
	ctx context.Context,
	binder framework.Binder,
	updateFn UpdateServiceBindingFunc,
	secretWriter SecretWriterFunc,
	getSvcBrokerFn refs.ServiceBrokerGetterFunc,
	getSvcClassFn refs.ServiceClassGetterFunc,
	getSvcInstanceFn refs.ServiceInstanceGetterFunc,
	evt watch.Event,
) error {

	serviceBinding := new(data.ServiceBinding)
	if err := data.TranslateToTPR(evt.Object, serviceBinding, data.ServiceBindingKind); err != nil {
		return ErrNotAServiceBinding
	}

	serviceBinding.Status.State = data.ServiceBindingStatePending
	serviceBinding, err := updateFn(serviceBinding)
	if err != nil {
		logger.Errorf("error updating service binding state to %s", serviceBinding.Status.State)
		return err
	}

	serviceBroker, serviceClass, serviceInstance, err := refs.GetDependenciesForServiceBinding(
		serviceBinding,
		getSvcBrokerFn,
		getSvcClassFn,
		getSvcInstanceFn,
	)
	if err != nil {
		logger.Errorf("getting service binding %s's dependencies (%s)", serviceBinding.Spec.ID, err)
		return err
	}

	bindReq := &framework.BindRequest{
		InstanceID: serviceInstance.Spec.ID,
		ServiceID:  serviceClass.ID,
		PlanID:     serviceInstance.Spec.PlanID,
		BindingID:  serviceBinding.Spec.ID,
		Parameters: serviceBinding.Spec.Parameters,
	}
	logger.Debugf("issuing bind request %+v", *bindReq)
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
			Name:      serviceBinding.Spec.SecretName,
			Namespace: serviceBinding.Namespace,
		},
		Data: bindResponseCredsToSecretData(bindResp.Creds),
	}

	if _, err := secretWriter(secret); err != nil {
		logger.Errorf("writing secret %s (%s)", serviceBinding.Spec.SecretName, err)
		return err
	}

	serviceBinding.Status.State = data.ServiceBindingStateBound
	if _, err := updateFn(serviceBinding); err != nil {
		logger.Errorf("error updating service binding state to %s", serviceBinding.Status.State)
		return err
	}

	return nil
}

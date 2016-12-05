package serviceinstance

import (
	"context"
	"errors"

	"github.com/deis/steward-framework"
	"github.com/deis/steward-framework/k8s/data"
	"github.com/deis/steward-framework/k8s/refs"
	"k8s.io/client-go/pkg/watch"
)

var (
	ErrCancelled           = errors.New("stopped")
	ErrNotAServiceInstance = errors.New("not an service instance")
	ErrWatchClosed         = errors.New("watch closed")
)

// RunLoop starts a blocking control loop that watches and takes action on service instance
// resources
func RunLoop(
	ctx context.Context,
	watchFn WatchServiceInstanceFunc,
	updateFn UpdateServiceInstanceFunc,
	getServiceClassFn refs.ServiceClassGetterFunc,
	getServiceBrokerFn refs.ServiceBrokerGetterFunc,
	lifecycler framework.Lifecycler,
) error {
	// TODO: We should watch ALL namespaces; this is temporary.
	// See https://github.com/deis/steward-framework/issues/29
	watcher, err := watchFn("default")
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
				logger.Errorf("service instance loop watch channel was closed")
				return ErrWatchClosed
			}
			logger.Debugf("service instance loop received event")
			switch evt.Type {
			case watch.Added:
				if err := handleAddServiceInstance(
					ctx,
					lifecycler,
					updateFn,
					getServiceClassFn,
					getServiceBrokerFn,
					evt,
				); err != nil {
					// TODO: try the handler again.
					// See https://github.com/deis/steward-framework/issues/26
					logger.Errorf("add service instance event handler failed (%s)", err)
				}
			case watch.Deleted:
				if err := handleDeleteServiceInstance(
					ctx,
					lifecycler,
					getServiceClassFn,
					getServiceBrokerFn,
					evt,
				); err != nil {
					// TODO: try the handler again.
					// See https://github.com/deis/steward-framework/issues/26
					logger.Errorf("delete service instance event handler failed (%s)", err)
				}
			}
		}
	}
}

func handleAddServiceInstance(
	ctx context.Context,
	lifecycler framework.Lifecycler,
	updateFn UpdateServiceInstanceFunc,
	getServiceClassFn refs.ServiceClassGetterFunc,
	getServiceBrokerFn refs.ServiceBrokerGetterFunc,
	evt watch.Event,
) error {

	serviceInstance := new(data.ServiceInstance)
	if err := data.TranslateToTPR(evt.Object, serviceInstance, data.ServiceInstanceKind); err != nil {
		return ErrNotAServiceInstance
	}
	serviceInstance.Status.Status = data.ServiceInstanceStatePending
	serviceInstance, err := updateFn(serviceInstance)
	if err != nil {
		return err
	}
	sc, err := getServiceClassFn(serviceInstance.Spec.ServiceClassRef)
	if err != nil {
		scNamespace := serviceInstance.Spec.ServiceClassRef.Namespace
		scName := serviceInstance.Spec.ServiceClassRef.Name
		logger.Errorf("couldn't find service class %s/%s", scNamespace, scName)
		return err
	}
	b, err := getServiceBrokerFn(sc.ServiceBrokerRef)
	if err != nil {
		serviceBrokerNamespace := sc.ServiceBrokerRef.Namespace
		serviceBrokerName := sc.ServiceBrokerRef.Name
		logger.Errorf("couldn't find service broker %s/%s", serviceBrokerNamespace, serviceBrokerName)
		return err
	}
	req := &framework.ProvisionRequest{
		InstanceID:        serviceInstance.Spec.ID,
		ServiceID:         sc.ID,
		PlanID:            serviceInstance.Spec.PlanID,
		AcceptsIncomplete: false,
		Parameters:        serviceInstance.Spec.Parameters,
	}
	finalServiceInstanceState := data.ServiceInstanceStateFailed
	defer func() {
		serviceInstance.Status.Status = finalServiceInstanceState
		if _, err = updateFn(serviceInstance); err != nil {
			logger.Errorf("failed to update service instance to final state %s (%s)", finalServiceInstanceState, err)
		}
	}()
	logger.Debugf("issuing provisioning request %+v", *req)
	_, err = lifecycler.Provision(ctx, b.Spec, req)
	if err != nil {
		return err
	}
	// TODO: Wait for async provision completion.
	// See https://github.com/deis/steward-framework/issues/39
	finalServiceInstanceState = data.ServiceInstanceStateProvisioned
	return nil
}

func handleDeleteServiceInstance(
	ctx context.Context,
	lifecycler framework.Lifecycler,
	getServiceClassFn refs.ServiceClassGetterFunc,
	getServiceBrokerFn refs.ServiceBrokerGetterFunc,
	evt watch.Event,
) error {
	serviceInstance, ok := evt.Object.(*data.ServiceInstance)
	if !ok {
		return ErrNotAServiceInstance
	}
	sc, err := getServiceClassFn(serviceInstance.Spec.ServiceClassRef)
	if err != nil {
		return err
	}
	b, err := getServiceBrokerFn(sc.ServiceBrokerRef)
	if err != nil {
		return err
	}
	req := &framework.DeprovisionRequest{
		InstanceID:        serviceInstance.Spec.ID,
		ServiceID:         sc.ID,
		PlanID:            serviceInstance.Spec.PlanID,
		AcceptsIncomplete: false,
		Parameters:        serviceInstance.Spec.Parameters,
	}
	_, err = lifecycler.Deprovision(ctx, b.Spec, req)
	if err != nil {
		return err
	}
	// TODO: Wait for async deprovision completion
	// See https://github.com/deis/steward-framework/issues/39
	return nil
}

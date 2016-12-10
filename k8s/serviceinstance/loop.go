package serviceinstance

import (
	"context"
	"errors"
	"time"

	"github.com/deis/steward-framework"
	"github.com/deis/steward-framework/k8s/data"
	"github.com/deis/steward-framework/k8s/refs"
	"k8s.io/client-go/pkg/watch"
)

var (
	ErrCancelled           = errors.New("stopped")
	ErrNotAServiceInstance = errors.New("not a service instance")
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
				watcher, err = watchFn("default")
				if err != nil {
					logger.Errorf("service instance loop watch channel was closed")
					return ErrWatchClosed
				}
				ch = watcher.ResultChan()
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
		AcceptsIncomplete: true,
		Parameters:        serviceInstance.Spec.Parameters,
	}
	finalServiceInstanceState := data.ServiceInstanceStateFailed
	finalServiceInstanceStateReason := ""
	defer func() {
		serviceInstance.Status.Status = finalServiceInstanceState
		serviceInstance.Status.StatusReason = finalServiceInstanceStateReason
		if _, err = updateFn(serviceInstance); err != nil {
			logger.Errorf("failed to update service instance to final state %s (%s)", finalServiceInstanceState, err)
		}
	}()
	logger.Debugf("issuing provisioning request %+v", *req)
	provisionResp, err := lifecycler.Provision(ctx, b.Spec, req)
	if err != nil {
		return err
	}
	if provisionResp.IsAsync {
		serviceInstance.Status.StatusReason = "Asynchronous provisioning in progress"
		if serviceInstance, err = updateFn(serviceInstance); err != nil {
			logger.Errorf("failed to update service instance (%s)", err)
			// Don't return the error if one occurs here. Just because we couldn't set set the status
			// reason here doesn't mean we shouldn't proceed with polling to attempt to resolve the final
			// state... after all, provisioning is on-going at this point. We should be making a best
			// effort to understand if it succeeds or not.
		}
		// TODO: Make this configurable
		ctxWithTimeout, cancelFn := context.WithTimeout(ctx, time.Hour)
		defer cancelFn()
		finalServiceInstanceState, finalServiceInstanceStateReason, err = pollProvisionState(ctxWithTimeout, b.Spec, lifecycler, serviceInstance.Spec.ID, sc.ID, serviceInstance.Spec.PlanID)
		if err != nil {
			return err
		}
	} else {
		finalServiceInstanceState = data.ServiceInstanceStateProvisioned
	}
	return nil
}

func handleDeleteServiceInstance(
	ctx context.Context,
	lifecycler framework.Lifecycler,
	getServiceClassFn refs.ServiceClassGetterFunc,
	getServiceBrokerFn refs.ServiceBrokerGetterFunc,
	evt watch.Event,
) error {
	serviceInstance := new(data.ServiceInstance)
	if err := data.TranslateToTPR(evt.Object, serviceInstance, data.ServiceInstanceKind); err != nil {
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
		AcceptsIncomplete: true,
		Parameters:        serviceInstance.Spec.Parameters,
	}
	deprovisionResp, err := lifecycler.Deprovision(ctx, b.Spec, req)
	if err != nil {
		return err
	}
	if deprovisionResp.IsAsync {
		// TODO: Make this configurable
		ctxWithTimeout, cancelFn := context.WithTimeout(ctx, time.Hour)
		defer cancelFn()
		err = pollDeprovisionState(ctxWithTimeout, b.Spec, lifecycler, serviceInstance.Spec.ID, sc.ID, serviceInstance.Spec.PlanID)
		if err != nil {
			return err
		}
	}
	return nil
}

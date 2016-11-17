package instance

import (
	"context"
	"errors"

	"github.com/deis/steward-framework"
	"github.com/deis/steward-framework/k8s/data"
	"github.com/deis/steward-framework/k8s/refs"
	"k8s.io/client-go/pkg/watch"
)

var (
	ErrCancelled     = errors.New("stopped")
	ErrNotAnInstance = errors.New("not an instance")
	ErrWatchClosed   = errors.New("watch closed")
)

// RunLoop starts a blocking control loop that watches and takes action on instance resources
func RunLoop(
	ctx context.Context,
	watchFn WatchInstanceFunc,
	updateFn UpdateInstanceFunc,
	getServiceClassFn refs.ServiceClassGetterFunc,
	getBrokerFn refs.BrokerGetterFunc,
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
				logger.Errorf("instance loop watch channel was closed")
				return ErrWatchClosed
			}
			logger.Debugf("instance loop received event")
			switch evt.Type {
			case watch.Added:
				if err := handleAddInstance(
					ctx,
					lifecycler,
					updateFn,
					getServiceClassFn,
					getBrokerFn,
					evt,
				); err != nil {
					// TODO: try the handler again.
					// See https://github.com/deis/steward-framework/issues/26
					logger.Errorf("add instance event handler failed (%s)", err)
				}
			case watch.Deleted:
				if err := handleDeleteInstance(
					ctx,
					lifecycler,
					getServiceClassFn,
					getBrokerFn,
					evt,
				); err != nil {
					// TODO: try the handler again.
					// See https://github.com/deis/steward-framework/issues/26
					logger.Errorf("delete instance event handler failed (%s)", err)
				}
			}
		}
	}
}

func handleAddInstance(
	ctx context.Context,
	lifecycler framework.Lifecycler,
	updateFn UpdateInstanceFunc,
	getServiceClassFn refs.ServiceClassGetterFunc,
	getBrokerFn refs.BrokerGetterFunc,
	evt watch.Event,
) error {
	instance, ok := evt.Object.(*data.Instance)
	if !ok {
		return ErrNotAnInstance
	}
	instance.Status.Status = data.InstanceStatePending
	if _, err := updateFn(instance); err != nil {
		return err
	}
	sc, err := getServiceClassFn(instance.Spec.ServiceClassRef)
	if err != nil {
		return err
	}
	b, err := getBrokerFn(sc.BrokerRef)
	if err != nil {
		return err
	}
	req := &framework.ProvisionRequest{
		InstanceID:        instance.Spec.ID,
		ServiceID:         sc.ID,
		PlanID:            instance.Spec.PlanID,
		AcceptsIncomplete: false,
		Parameters:        instance.Spec.Parameters,
	}
	finalInstanceState := data.InstanceStateFailed
	defer func() {
		instance.Status.Status = finalInstanceState
		if _, err = updateFn(instance); err != nil {
			logger.Errorf("failed to update instance to state %s (%s)", finalInstanceState, err)
		}
	}()
	_, err = lifecycler.Provision(ctx, b.Spec, req)
	if err != nil {
		return err
	}
	// TODO: Wait for async provision completion.
	// See https://github.com/deis/steward-framework/issues/39
	finalInstanceState = data.InstanceStateProvisioned
	return nil
}

func handleDeleteInstance(
	ctx context.Context,
	lifecycler framework.Lifecycler,
	getServiceClassFn refs.ServiceClassGetterFunc,
	getBrokerFn refs.BrokerGetterFunc,
	evt watch.Event,
) error {
	instance, ok := evt.Object.(*data.Instance)
	if !ok {
		return ErrNotAnInstance
	}
	sc, err := getServiceClassFn(instance.Spec.ServiceClassRef)
	if err != nil {
		return err
	}
	b, err := getBrokerFn(sc.BrokerRef)
	if err != nil {
		return err
	}
	req := &framework.DeprovisionRequest{
		InstanceID:        instance.Spec.ID,
		ServiceID:         sc.ID,
		PlanID:            instance.Spec.PlanID,
		AcceptsIncomplete: false,
		Parameters:        instance.Spec.Parameters,
	}
	_, err = lifecycler.Deprovision(ctx, b.Spec, req)
	if err != nil {
		return err
	}
	// TODO: Wait for async deprovision completion
	// See https://github.com/deis/steward-framework/issues/39
	return nil
}

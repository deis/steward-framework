package serviceinstance

import (
	"context"
	"testing"
	"time"

	"github.com/arschles/assert"
	"github.com/deis/steward-framework/fake"
	"github.com/deis/steward-framework/k8s/data"
	"github.com/deis/steward-framework/k8s/refs"
	"k8s.io/client-go/pkg/api"
	"k8s.io/client-go/pkg/watch"
)

func TestRunLoopCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	watcher, fakeWatcher := newFakeWatchServiceInstanceFunc(nil)
	updater, updated := newFakeUpdateServiceInstanceFunc(nil)
	getServiceClassFn := refs.NewFakeServiceClassGetterFunc(&data.ServiceClass{}, nil)
	getServiceBrokerFn := refs.NewFakeServiceBrokerGetterFunc(&data.ServiceBroker{}, nil)
	lifecycler := &fake.Lifecycler{}
	assert.Err(t, ErrCancelled, RunLoop(ctx, watcher, updater, getServiceClassFn, getServiceBrokerFn, lifecycler))
	assert.True(t, fakeWatcher.IsStopped(), "watcher isn't stopped")
	assert.Equal(t, len(*updated), 0, "number of updated service instances")
}

func TestRunLoopAddSuccess(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	watcher, fakeWatcher := newFakeWatchServiceInstanceFunc(nil)
	updater, updated := newFakeUpdateServiceInstanceFunc(nil)
	getServiceClassFn := refs.NewFakeServiceClassGetterFunc(&data.ServiceClass{}, nil)
	getServiceBrokerFn := refs.NewFakeServiceBrokerGetterFunc(&data.ServiceBroker{}, nil)
	provisioner := fake.NewProvisioner()
	lifecycler := &fake.Lifecycler{
		Provisioner: provisioner,
	}

	errCh := make(chan error)
	go func() {
		errCh <- RunLoop(ctx, watcher, updater, getServiceClassFn, getServiceBrokerFn, lifecycler)
	}()

	inst := new(data.ServiceInstance)
	inst.Kind = data.ServiceInstanceKind
	fakeWatcher.Add(inst)

	time.Sleep(100 * time.Millisecond)
	cancel()
	fakeWatcher.Stop()

	err := <-errCh
	assert.Equal(t, len(provisioner.Reqs), 1, "number of provision requests")
	assert.Equal(t, len(*updated), 2, "number of updated service instances")
	assert.Err(t, ErrCancelled, err)
}

func TestRunLoopDeleteSuccess(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	watcher, fakeWatcher := newFakeWatchServiceInstanceFunc(nil)
	updater, updated := newFakeUpdateServiceInstanceFunc(nil)
	getServiceClassFn := refs.NewFakeServiceClassGetterFunc(&data.ServiceClass{}, nil)
	getServiceBrokerFn := refs.NewFakeServiceBrokerGetterFunc(&data.ServiceBroker{}, nil)
	deprovisioner := fake.NewDeprovisioner()
	lifecycler := &fake.Lifecycler{
		Deprovisioner: deprovisioner,
	}

	errCh := make(chan error)
	go func() {
		errCh <- RunLoop(ctx, watcher, updater, getServiceClassFn, getServiceBrokerFn, lifecycler)
	}()

	inst := new(data.ServiceInstance)
	inst.Kind = data.ServiceInstanceKind
	fakeWatcher.Delete(inst)
	time.Sleep(100 * time.Millisecond)
	cancel()
	fakeWatcher.Stop()

	err := <-errCh
	assert.Equal(t, len(deprovisioner.Reqs), 1, "number of deprovision requests")
	assert.Equal(t, len(*updated), 0, "number of updated service instances")
	assert.Err(t, ErrCancelled, err)
}

func TestHandleAddServiceInstanceNotAServiceInstance(t *testing.T) {
	ctx := context.Background()
	updater, updated := newFakeUpdateServiceInstanceFunc(nil)
	getServiceClassFn := refs.NewFakeServiceClassGetterFunc(&data.ServiceClass{}, nil)
	getServiceBrokerFn := refs.NewFakeServiceBrokerGetterFunc(&data.ServiceBroker{}, nil)
	provisioner := fake.NewProvisioner()
	lifecycler := &fake.Lifecycler{
		Provisioner: provisioner,
	}
	evt := watch.Event{
		Type:   watch.Added,
		Object: &api.Pod{},
	}
	err := handleAddServiceInstance(ctx, lifecycler, updater, getServiceClassFn, getServiceBrokerFn, evt)
	assert.Err(t, ErrNotAServiceInstance, err)
	assert.Equal(t, len(provisioner.Reqs), 0, "number of provision requests")
	assert.Equal(t, len(*updated), 0, "number of updated service instances")
}

func TestHandleAddServiceInstanceSuccess(t *testing.T) {
	ctx := context.Background()
	updater, updated := newFakeUpdateServiceInstanceFunc(nil)
	getServiceClassFn := refs.NewFakeServiceClassGetterFunc(&data.ServiceClass{}, nil)
	getServiceBrokerFn := refs.NewFakeServiceBrokerGetterFunc(&data.ServiceBroker{}, nil)
	provisioner := fake.NewProvisioner()
	lifecycler := &fake.Lifecycler{
		Provisioner: provisioner,
	}
	inst := new(data.ServiceInstance)
	inst.Kind = data.ServiceInstanceKind
	evt := watch.Event{
		Type:   watch.Added,
		Object: inst,
	}
	err := handleAddServiceInstance(ctx, lifecycler, updater, getServiceClassFn, getServiceBrokerFn, evt)
	assert.NoErr(t, err)
	assert.Equal(t, len(provisioner.Reqs), 1, "number of provision requests")
	assert.Equal(t, len(*updated), 2, "number of updated service instances")
}

func TestHandleDeleteServiceInstanceNotAServiceInstance(t *testing.T) {
	ctx := context.Background()
	getServiceClassFn := refs.NewFakeServiceClassGetterFunc(&data.ServiceClass{}, nil)
	getServiceBrokerFn := refs.NewFakeServiceBrokerGetterFunc(&data.ServiceBroker{}, nil)
	deprovisioner := fake.NewDeprovisioner()
	lifecycler := &fake.Lifecycler{
		Deprovisioner: deprovisioner,
	}
	evt := watch.Event{
		Type:   watch.Deleted,
		Object: &api.Pod{},
	}
	err := handleDeleteServiceInstance(ctx, lifecycler, getServiceClassFn, getServiceBrokerFn, evt)
	assert.Err(t, ErrNotAServiceInstance, err)
	assert.Equal(t, len(deprovisioner.Reqs), 0, "number of deprovision requests")
}

func TestHandleDeleteServiceInstanceSuccess(t *testing.T) {
	ctx := context.Background()
	getServiceClassFn := refs.NewFakeServiceClassGetterFunc(&data.ServiceClass{}, nil)
	getServiceBrokerFn := refs.NewFakeServiceBrokerGetterFunc(&data.ServiceBroker{}, nil)
	deprovisioner := fake.NewDeprovisioner()
	lifecycler := &fake.Lifecycler{
		Deprovisioner: deprovisioner,
	}
	inst := new(data.ServiceInstance)
	inst.Kind = data.ServiceInstanceKind
	evt := watch.Event{
		Type:   watch.Deleted,
		Object: inst,
	}
	err := handleDeleteServiceInstance(ctx, lifecycler, getServiceClassFn, getServiceBrokerFn, evt)
	assert.NoErr(t, err)
	assert.Equal(t, len(deprovisioner.Reqs), 1, "number of deprovision requests")
}

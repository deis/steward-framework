package instance

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
	watcher, fakeWatcher := newFakeWatchInstanceFunc(nil)
	updater, updated := newFakeUpdateInstanceFunc(nil)
	getServiceClassFn := refs.NewFakeServiceClassGetterFunc(&data.ServiceClass{}, nil)
	getServiceBrokerFn := refs.NewFakeServiceBrokerGetterFunc(&data.ServiceBroker{}, nil)
	lifecycler := &fake.Lifecycler{}
	assert.Err(t, ErrCancelled, RunLoop(ctx, watcher, updater, getServiceClassFn, getServiceBrokerFn, lifecycler))
	assert.True(t, fakeWatcher.IsStopped(), "watcher isn't stopped")
	assert.Equal(t, len(*updated), 0, "number of updated instances")
}

func TestRunLoopAddSuccess(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	watcher, fakeWatcher := newFakeWatchInstanceFunc(nil)
	updater, updated := newFakeUpdateInstanceFunc(nil)
	getServiceClassFn := refs.NewFakeServiceClassGetterFunc(&data.ServiceClass{}, nil)
	getServiceBrokerFn := refs.NewFakeServiceBrokerGetterFunc(&data.ServiceBroker{}, nil)
	provisioner := &fake.Provisioner{}
	lifecycler := &fake.Lifecycler{
		Provisioner: provisioner,
	}

	errCh := make(chan error)
	go func() {
		errCh <- RunLoop(ctx, watcher, updater, getServiceClassFn, getServiceBrokerFn, lifecycler)
	}()

	inst := new(data.Instance)
	inst.Kind = data.InstanceKind
	fakeWatcher.Add(inst)
	fakeWatcher.Stop()

	const errTimeout = 100 * time.Millisecond
	select {
	case err := <-errCh:
		assert.Equal(t, len(provisioner.Reqs), 1, "number of provision requests")
		assert.Equal(t, len(*updated), 2, "number of updated instances")
		assert.Err(t, ErrWatchClosed, err)
	case <-time.After(errTimeout):
		t.Fatalf("RunLoop didn't return after %s", errTimeout)
	}

}

func TestRunLoopDeleteSuccess(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	watcher, fakeWatcher := newFakeWatchInstanceFunc(nil)
	updater, updated := newFakeUpdateInstanceFunc(nil)
	getServiceClassFn := refs.NewFakeServiceClassGetterFunc(&data.ServiceClass{}, nil)
	getServiceBrokerFn := refs.NewFakeServiceBrokerGetterFunc(&data.ServiceBroker{}, nil)
	deprovisioner := &fake.Deprovisioner{}
	lifecycler := &fake.Lifecycler{
		Deprovisioner: deprovisioner,
	}

	errCh := make(chan error)
	go func() {
		errCh <- RunLoop(ctx, watcher, updater, getServiceClassFn, getServiceBrokerFn, lifecycler)
	}()

	fakeWatcher.Delete(&data.Instance{})
	fakeWatcher.Stop()

	const errTimeout = 100 * time.Millisecond
	select {
	case err := <-errCh:
		assert.Equal(t, len(deprovisioner.Reqs), 1, "number of deprovision requests")
		assert.Equal(t, len(*updated), 0, "number of updated instances")
		assert.Err(t, ErrWatchClosed, err)
	case <-time.After(errTimeout):
		t.Fatalf("RunLoop didn't return after %s", errTimeout)
	}

}

func TestHandleAddInstanceNotAnInstance(t *testing.T) {
	ctx := context.Background()
	updater, updated := newFakeUpdateInstanceFunc(nil)
	getServiceClassFn := refs.NewFakeServiceClassGetterFunc(&data.ServiceClass{}, nil)
	getServiceBrokerFn := refs.NewFakeServiceBrokerGetterFunc(&data.ServiceBroker{}, nil)
	provisioner := &fake.Provisioner{}
	lifecycler := &fake.Lifecycler{
		Provisioner: provisioner,
	}
	evt := watch.Event{
		Type:   watch.Added,
		Object: &api.Pod{},
	}
	err := handleAddInstance(ctx, lifecycler, updater, getServiceClassFn, getServiceBrokerFn, evt)
	assert.Err(t, ErrNotAnInstance, err)
	assert.Equal(t, len(provisioner.Reqs), 0, "number of provision requests")
	assert.Equal(t, len(*updated), 0, "number of updated instances")
}

func TestHandleAddInstanceSuccess(t *testing.T) {
	ctx := context.Background()
	updater, updated := newFakeUpdateInstanceFunc(nil)
	getServiceClassFn := refs.NewFakeServiceClassGetterFunc(&data.ServiceClass{}, nil)
	getServiceBrokerFn := refs.NewFakeServiceBrokerGetterFunc(&data.ServiceBroker{}, nil)
	provisioner := &fake.Provisioner{}
	lifecycler := &fake.Lifecycler{
		Provisioner: provisioner,
	}
	inst := new(data.Instance)
	inst.Kind = data.InstanceKind
	evt := watch.Event{
		Type:   watch.Added,
		Object: inst,
	}
	err := handleAddInstance(ctx, lifecycler, updater, getServiceClassFn, getServiceBrokerFn, evt)
	assert.NoErr(t, err)
	assert.Equal(t, len(provisioner.Reqs), 1, "number of provision requests")
	assert.Equal(t, len(*updated), 2, "number of updated instances")
}

func TestHandleDeleteInstanceNotAnInstance(t *testing.T) {
	ctx := context.Background()
	getServiceClassFn := refs.NewFakeServiceClassGetterFunc(&data.ServiceClass{}, nil)
	getServiceBrokerFn := refs.NewFakeServiceBrokerGetterFunc(&data.ServiceBroker{}, nil)
	deprovisioner := &fake.Deprovisioner{}
	lifecycler := &fake.Lifecycler{
		Deprovisioner: deprovisioner,
	}
	evt := watch.Event{
		Type:   watch.Deleted,
		Object: &api.Pod{},
	}
	err := handleDeleteInstance(ctx, lifecycler, getServiceClassFn, getServiceBrokerFn, evt)
	assert.Err(t, ErrNotAnInstance, err)
	assert.Equal(t, len(deprovisioner.Reqs), 0, "number of deprovision requests")
}

func TestHandleDeleteInstanceSuccess(t *testing.T) {
	ctx := context.Background()
	getServiceClassFn := refs.NewFakeServiceClassGetterFunc(&data.ServiceClass{}, nil)
	getServiceBrokerFn := refs.NewFakeServiceBrokerGetterFunc(&data.ServiceBroker{}, nil)
	deprovisioner := &fake.Deprovisioner{}
	lifecycler := &fake.Lifecycler{
		Deprovisioner: deprovisioner,
	}
	evt := watch.Event{
		Type:   watch.Deleted,
		Object: &data.Instance{},
	}
	err := handleDeleteInstance(ctx, lifecycler, getServiceClassFn, getServiceBrokerFn, evt)
	assert.NoErr(t, err)
	assert.Equal(t, len(deprovisioner.Reqs), 1, "number of deprovision requests")
}

package servicebinding

import (
	"context"
	"testing"
	"time"

	"github.com/arschles/assert"
	"github.com/deis/steward-framework"
	"github.com/deis/steward-framework/fake"
	"github.com/deis/steward-framework/k8s/data"
	"github.com/deis/steward-framework/k8s/refs"
	"k8s.io/client-go/pkg/api"
	"k8s.io/client-go/pkg/api/unversioned"
	"k8s.io/client-go/pkg/watch"
)

func TestRunLoopCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	watcher, fakeWatcher := newFakeWatchServiceBindingFunc(nil)
	secretWriter, writtenSecrets := newFakeSecretWriterFunc(nil)
	updateServiceBindingFn, updatedServiceBindings := newFakeUpdateServiceBindingFunc(nil)
	binder := &fake.Binder{}
	assert.Err(t, ErrCancelled, RunLoop(
		ctx,
		"testns",
		binder,
		secretWriter,
		watcher,
		updateServiceBindingFn,
		refs.NewFakeServiceBrokerGetterFunc(nil, nil),
		refs.NewFakeServiceClassGetterFunc(nil, nil),
		refs.NewFakeServiceInstanceGetterFunc(nil, nil),
	))
	assert.Equal(t, len(binder.Reqs), 0, "number of bind calls")
	assert.True(t, fakeWatcher.IsStopped(), "watcher was not stopped")
	assert.Equal(t, len(*writtenSecrets), 0, "number of secrets written")
	assert.Equal(t, len(*updatedServiceBindings), 0, "number of service bindings written")
}

func TestRunLoopSuccess(t *testing.T) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	watcher, fakeWatcher := newFakeWatchServiceBindingFunc(nil)
	secretWriter, writtenSecrets := newFakeSecretWriterFunc(nil)
	updateServiceBindingFn, updatedServiceBindings := newFakeUpdateServiceBindingFunc(nil)
	binder := &fake.Binder{
		Res: &framework.BindResponse{
			Creds: map[string]interface{}{"a": "b"},
		},
	}
	svcBrokerGetterFn := refs.NewFakeServiceBrokerGetterFunc(&data.ServiceBroker{}, nil)
	svcClassGetterFn := refs.NewFakeServiceClassGetterFunc(&data.ServiceClass{}, nil)
	svcInstanceGetterFn := refs.NewFakeServiceInstanceGetterFunc(&data.ServiceInstance{}, nil)

	errCh := make(chan error)
	go func() {
		errCh <- RunLoop(
			ctx,
			"testns",
			binder,
			secretWriter,
			watcher,
			updateServiceBindingFn,
			svcBrokerGetterFn,
			svcClassGetterFn,
			svcInstanceGetterFn,
		)
	}()
	serviceBinding := new(data.ServiceBinding)
	serviceBinding.Kind = data.ServiceBindingKind
	fakeWatcher.Add(serviceBinding)
	fakeWatcher.Stop()

	const errTimeout = 100 * time.Millisecond
	select {
	case err := <-errCh:
		assert.Equal(t, len(*writtenSecrets), 1, "number of written secrets")
		writtenSecret := (*writtenSecrets)[0]
		assert.Equal(t, len(writtenSecret.Data), len(binder.Res.Creds), "number of credentials in the secret")
		assert.Equal(t, len(*updatedServiceBindings), 2, "number of updated service bindings")
		assert.Equal(t, len(binder.Reqs), 1, "number of bind calls")
		assert.Err(t, ErrWatchClosed, err)
	case <-time.After(errTimeout):
		t.Fatalf("RunLoop didn't return after %s", errTimeout)
	}
}

func TestHandleAddServiceBindingNotAServiceBinding(t *testing.T) {
	evt := watch.Event{
		Type: watch.Added,
		Object: &api.Pod{
			TypeMeta: unversioned.TypeMeta{Kind: "Pod"},
		},
	}
	err := handleAddServiceBinding(
		context.Background(),
		nil, // framework.Binder
		nil, // UpdateServiceBindingFunc,
		nil, // SecretWriterFunc
		nil, // refs.ServiceBrokerGetterFunc
		nil, // refs.ServiceClassGetterFunc
		nil, // refs.ServiceInstanceGetterFunc
		evt,
	)
	assert.Err(t, ErrNotAServiceBinding, err)
}

func TestHandleAddServiceBindingSuccess(t *testing.T) {
	ctx := context.Background()
	binder := &fake.Binder{Res: &framework.BindResponse{}}
	updateServiceBindingFn, updatedServiceBindings := newFakeUpdateServiceBindingFunc(nil)
	secretWriterFn, writtenSecrets := newFakeSecretWriterFunc(nil)
	serviceBrokerGetterFn := refs.NewFakeServiceBrokerGetterFunc(&data.ServiceBroker{}, nil)
	svcClassGetterFn := refs.NewFakeServiceClassGetterFunc(&data.ServiceClass{}, nil)
	svcInstanceGetterFn := refs.NewFakeServiceInstanceGetterFunc(&data.ServiceInstance{}, nil)
	serviceBinding := new(data.ServiceBinding)
	serviceBinding.Kind = data.ServiceBindingKind
	evt := watch.Event{
		Type:   watch.Added,
		Object: serviceBinding,
	}

	err := handleAddServiceBinding(
		ctx,
		binder,
		updateServiceBindingFn,
		secretWriterFn,
		serviceBrokerGetterFn,
		svcClassGetterFn,
		svcInstanceGetterFn,
		evt,
	)

	assert.NoErr(t, err)
	assert.Equal(t, len(binder.Reqs), 1, "number of bind requests")
	assert.Equal(t, len(*updatedServiceBindings), 2, "number of updated service binding resources")
	assert.Equal(t, len(*writtenSecrets), 1, "number of written secrets")
}

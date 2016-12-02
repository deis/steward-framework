package binding

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
	watcher, fakeWatcher := newFakeWatchBindingFunc(nil)
	secretWriter, writtenSecrets := newFakeSecretWriterFunc(nil)
	updateBindingFn, updatedBindings := newFakeUpdateBindingFunc(nil)
	binder := &fake.Binder{}
	assert.Err(t, ErrCancelled, RunLoop(
		ctx,
		"testns",
		binder,
		secretWriter,
		watcher,
		updateBindingFn,
		refs.NewFakeServiceBrokerGetterFunc(nil, nil),
		refs.NewFakeServiceClassGetterFunc(nil, nil),
		refs.NewFakeInstanceGetterFunc(nil, nil),
	))
	assert.Equal(t, len(binder.Reqs), 0, "number of bind calls")
	assert.True(t, fakeWatcher.IsStopped(), "watcher was not stopped")
	assert.Equal(t, len(*writtenSecrets), 0, "number of secrets written")
	assert.Equal(t, len(*updatedBindings), 0, "number of bindings written")
}

func TestRunLoopSuccess(t *testing.T) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	watcher, fakeWatcher := newFakeWatchBindingFunc(nil)
	secretWriter, writtenSecrets := newFakeSecretWriterFunc(nil)
	updateBindingFn, updatedBindings := newFakeUpdateBindingFunc(nil)
	binder := &fake.Binder{
		Res: &framework.BindResponse{
			Creds: map[string]interface{}{"a": "b"},
		},
	}
	svcBrokerGetterFn := refs.NewFakeServiceBrokerGetterFunc(&data.ServiceBroker{}, nil)
	svcClassGetterFn := refs.NewFakeServiceClassGetterFunc(&data.ServiceClass{}, nil)
	instanceGetterFn := refs.NewFakeInstanceGetterFunc(&data.Instance{}, nil)

	errCh := make(chan error)
	go func() {
		errCh <- RunLoop(
			ctx,
			"testns",
			binder,
			secretWriter,
			watcher,
			updateBindingFn,
			svcBrokerGetterFn,
			svcClassGetterFn,
			instanceGetterFn,
		)
	}()
	binding := new(data.Binding)
	binding.Kind = data.BindingKind
	fakeWatcher.Add(binding)
	fakeWatcher.Stop()

	const errTimeout = 100 * time.Millisecond
	select {
	case err := <-errCh:
		assert.Equal(t, len(*writtenSecrets), 1, "number of written secrets")
		writtenSecret := (*writtenSecrets)[0]
		assert.Equal(t, len(writtenSecret.Data), len(binder.Res.Creds), "number of credentials in the secret")
		assert.Equal(t, len(*updatedBindings), 2, "number of updated bindings")
		assert.Equal(t, len(binder.Reqs), 1, "number of bind calls")
		assert.Err(t, ErrWatchClosed, err)
	case <-time.After(errTimeout):
		t.Fatalf("RunLoop didn't return after %s", errTimeout)
	}
}

func TestHandleAddBindingNotABinding(t *testing.T) {
	evt := watch.Event{
		Type: watch.Added,
		Object: &api.Pod{
			TypeMeta: unversioned.TypeMeta{Kind: "Pod"},
		},
	}
	err := handleAddBinding(
		context.Background(),
		nil, // framework.Binder
		nil, // UpdateBindingFunc,
		nil, // SecretWriterFunc
		nil, // refs.ServiceBrokerGetterFunc
		nil, // refs.ServiceClassGetterFunc
		nil, // refs.InstanceGetterFunc
		evt,
	)
	assert.Err(t, ErrNotABinding, err)
}

func TestHandleAddBindingSuccess(t *testing.T) {
	ctx := context.Background()
	binder := &fake.Binder{Res: &framework.BindResponse{}}
	updateBindingFn, updatedBindings := newFakeUpdateBindingFunc(nil)
	secretWriterFn, writtenSecrets := newFakeSecretWriterFunc(nil)
	serviceBrokerGetterFn := refs.NewFakeServiceBrokerGetterFunc(&data.ServiceBroker{}, nil)
	svcClassGetterFn := refs.NewFakeServiceClassGetterFunc(&data.ServiceClass{}, nil)
	instanceGetterFn := refs.NewFakeInstanceGetterFunc(&data.Instance{}, nil)
	binding := new(data.Binding)
	binding.Kind = data.BindingKind
	evt := watch.Event{
		Type:   watch.Added,
		Object: binding,
	}

	err := handleAddBinding(
		ctx,
		binder,
		updateBindingFn,
		secretWriterFn,
		serviceBrokerGetterFn,
		svcClassGetterFn,
		instanceGetterFn,
		evt,
	)

	assert.NoErr(t, err)
	assert.Equal(t, len(binder.Reqs), 1, "number of bind requests")
	assert.Equal(t, len(*updatedBindings), 2, "number of update Binding resources")
	assert.Equal(t, len(*writtenSecrets), 1, "number of written secrets")
}

package broker

import (
	"context"
	"testing"
	"time"

	"github.com/arschles/assert"
	"github.com/deis/steward-framework"
	"github.com/deis/steward-framework/fake"
	"github.com/deis/steward-framework/k8s/data"
	"k8s.io/client-go/pkg/api"
	"k8s.io/client-go/pkg/watch"
)

func fakeCataloger() fake.Cataloger {
	return fake.Cataloger{
		Services: []*framework.Service{
			&framework.Service{
				ServiceInfo: framework.ServiceInfo{
					Name:          "testName",
					ID:            "testID",
					Description:   "testDescr",
					PlanUpdatable: true,
				},
				Plans: []framework.ServicePlan{
					framework.ServicePlan{ID: "tid1", Name: "tName1", Description: "tDesc1", Free: true},
					framework.ServicePlan{ID: "tid2", Name: "tName2", Description: "tDesc2", Free: false},
					framework.ServicePlan{ID: "tid3", Name: "tName3", Description: "tDesc3", Free: true},
				},
			},
		},
	}
}

func TestRunLoopCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	watcher, fakeWatcher := newFakeWatchBrokerFunc(nil)
	assert.Err(t, ErrCancelled, RunLoop(ctx, "testns", watcher, nil, nil))
	assert.True(t, fakeWatcher.IsStopped(), "watcher isn't stopped")
}

func TestRunLoopSuccess(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	watcher, fakeWatcher := newFakeWatchBrokerFunc(nil)
	cataloger := fakeCataloger()
	createSvcClass, createdSvcClasses := newFakeCreateServiceClassFunc(nil)

	errCh := make(chan error)
	go func() {
		errCh <- RunLoop(ctx, "testns", watcher, cataloger, createSvcClass)
	}()

	fakeWatcher.Add(&data.Broker{})
	fakeWatcher.Stop()

	const errTimeout = 100 * time.Millisecond
	select {
	case err := <-errCh:
		assert.Equal(t, len(*createdSvcClasses), len(cataloger.Services), "number of created service classes")
		assert.Err(t, ErrWatchClosed, err)
	case <-time.After(errTimeout):
		t.Fatalf("RunLoop didn't return after %s", errTimeout)
	}

}

func TestHandleAddBrokerNotABroker(t *testing.T) {
	ctx := context.Background()
	cataloger := fake.Cataloger{}
	createSvcClass, createdSvcClasses := newFakeCreateServiceClassFunc(nil)
	evt := watch.Event{
		Type:   watch.Added,
		Object: &api.Pod{},
	}
	err := handleAddBroker(ctx, cataloger, createSvcClass, evt)
	assert.Err(t, ErrNotABroker, err)
	assert.Equal(t, len(*createdSvcClasses), 0, "number of create svc classes")
}

func TestHandleAddBrokerSuccess(t *testing.T) {
	ctx := context.Background()
	cataloger := fakeCataloger()
	createSvcClass, createdSvcClasses := newFakeCreateServiceClassFunc(nil)
	evt := watch.Event{
		Type:   watch.Added,
		Object: &data.Broker{},
	}
	err := handleAddBroker(ctx, cataloger, createSvcClass, evt)
	assert.NoErr(t, err)
	assert.Equal(t, len(*createdSvcClasses), len(cataloger.Services), "number of create svc classes")
}

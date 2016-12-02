package servicebroker

import (
	"context"
	"testing"
	"time"

	"github.com/arschles/assert"
	"github.com/deis/steward-framework"
	"github.com/deis/steward-framework/fake"
	"github.com/deis/steward-framework/k8s/data"
	"k8s.io/client-go/pkg/api"
	"k8s.io/client-go/pkg/api/unversioned"
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
	watcher, fakeWatcher := newFakeWatchServiceBrokerFunc(nil)
	updater, updated := newFakeUpdateServiceBrokerFunc(nil)
	assert.Err(t, ErrCancelled, RunLoop(ctx, "testns", watcher, updater, nil, nil))
	assert.True(t, fakeWatcher.IsStopped(), "watcher isn't stopped")
	assert.Equal(t, len(*updated), 0, "number of updated service brokers")
}

func TestRunLoopSuccess(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	watcher, fakeWatcher := newFakeWatchServiceBrokerFunc(nil)
	updater, updated := newFakeUpdateServiceBrokerFunc(nil)
	cataloger := fakeCataloger()
	createSvcClass, createdSvcClasses := newFakeCreateServiceClassFunc(nil)

	errCh := make(chan error)
	go func() {
		errCh <- RunLoop(ctx, "testns", watcher, updater, cataloger, createSvcClass)
	}()

	serviceBroker := new(data.ServiceBroker)
	serviceBroker.Kind = data.ServiceBrokerKind
	fakeWatcher.Add(serviceBroker)
	fakeWatcher.Stop()

	const errTimeout = 100 * time.Millisecond
	select {
	case err := <-errCh:
		assert.Equal(t, len(*createdSvcClasses), len(cataloger.Services), "number of created service classes")
		assert.Equal(t, len(*updated), 2, "number of updated service brokers")
		assert.Err(t, ErrWatchClosed, err)
	case <-time.After(errTimeout):
		t.Fatalf("RunLoop didn't return after %s", errTimeout)
	}

}

func TestHandleAddServiceBrokerNotAServiceBroker(t *testing.T) {
	ctx := context.Background()
	cataloger := fake.Cataloger{}
	updater, updated := newFakeUpdateServiceBrokerFunc(nil)
	createSvcClass, createdSvcClasses := newFakeCreateServiceClassFunc(nil)
	evt := watch.Event{
		Type: watch.Added,
		Object: &api.Pod{
			TypeMeta: unversioned.TypeMeta{Kind: "Pod"},
		},
	}
	err := handleAddServiceBroker(ctx, cataloger, updater, createSvcClass, evt)
	assert.Err(t, ErrNotAServiceBroker, err)
	assert.Equal(t, len(*createdSvcClasses), 0, "number of create service classes")
	assert.Equal(t, len(*updated), 0, "number of updated service brokers")
}

func TestHandleAddServiceBrokerSuccess(t *testing.T) {
	ctx := context.Background()
	cataloger := fakeCataloger()
	updater, updated := newFakeUpdateServiceBrokerFunc(nil)
	createSvcClass, createdSvcClasses := newFakeCreateServiceClassFunc(nil)
	serviceBroker := new(data.ServiceBroker)
	serviceBroker.Kind = data.ServiceBrokerKind
	evt := watch.Event{
		Type:   watch.Added,
		Object: serviceBroker,
	}
	err := handleAddServiceBroker(ctx, cataloger, updater, createSvcClass, evt)
	assert.NoErr(t, err)
	assert.Equal(t, len(*createdSvcClasses), len(cataloger.Services), "number of create svc classes")
	assert.Equal(t, len(*updated), 2, "number of updated service brokers")
}

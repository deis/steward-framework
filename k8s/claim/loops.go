package claim

import (
	"context"
	"time"

	"github.com/deis/steward-framework"
	"github.com/deis/steward-framework/k8s"
	"k8s.io/client-go/1.4/kubernetes/typed/core/v1"
)

var maxAsyncDuration = time.Hour

// StartControlLoops calls StartControlLoop for each namespace in namespaces. For each call to
// StartControlLoop, it calls evtNamespacer.Events(ns) to create a new ConfigMapInterface. Because
// each StartControlLoop call is done in a new goroutine, this function need not be called in its
// own goroutine.
//
// In order to stop all loops, pass a cancellable context to this function and call its cancel
// function
func StartControlLoops(
	ctx context.Context,
	evtNamespacer InteractorNamespacer,
	secretsNamespacer v1.SecretsGetter,
	lookup k8s.ServiceCatalogLookup,
	lifecycler framework.Lifecycler,
	namespaces []string,
	maximumAsyncDuration time.Duration,
	errCh chan<- error,
) {
	// krancour: Polling functions for checking async provisioning / deprovisioning status need to
	// know the maxAsyncDuration... however, the event loops we're about to start always trigger
	// functions of type nextFunc in response to events. I don't want to modify the nextFunc
	// signature to include arguments not used by most nextFuncs, so as a workaround,
	// maxAsyncDuration is copied here into a package-scoped variable where it will be available to
	// the polling functions when they need it.
	maxAsyncDuration = maximumAsyncDuration

	for _, ns := range namespaces {
		logger.Debugf("starting claims control loop for namespace %s", ns)
		go func(ns string) {
			evtIface := evtNamespacer.Interactor(ns)
			if err := StartControlLoop(ctx, evtIface, secretsNamespacer, lookup, lifecycler); err != nil {
				logger.Errorf("failed to start control loop for namespace %s, skipping (%s)", ns, err)
				errCh <- err
				return
			}
		}(ns)
	}
}

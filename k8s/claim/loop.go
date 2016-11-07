package claim

import (
	"context"
	"errors"

	"github.com/deis/steward-framework"
	"github.com/deis/steward-framework/k8s"
	"github.com/deis/steward-framework/k8s/claim/state"
	"k8s.io/client-go/kubernetes/typed/core/v1"
	v1types "k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/pkg/labels"
	"k8s.io/client-go/pkg/watch"
)

const (
	labelKeyType             = "type"
	labelValServicePlanClaim = "service-plan-claim"
)

var (
	claimLabelSelector = labels.SelectorFromSet(labels.Set(map[string]string{
		labelKeyType: labelValServicePlanClaim,
	}))
	errLoopStopped = errors.New("loop has been stopped")
	errWatchClosed = errors.New("watch closed")
)

// StartControlLoop starts an infinite loop that receives on watcher.ResultChan() and takes action
// on each change in service plan claims. It's intended to be called in a goroutine. Call
// watcher.Stop() to stop this loop
func StartControlLoop(
	ctx context.Context,
	iface Interactor,
	secretsNamespacer v1.SecretsGetter,
	lookup k8s.ServiceCatalogLookup,
	lifecycler framework.Lifecycler,
) error {

	// start up the watcher so that events build up on the channel while we're listing events
	// (which happens below)
	cancelCtx, cancelFn := context.WithCancel(ctx)
	watcher := iface.Watch(cancelCtx, v1types.ListOptions{LabelSelector: claimLabelSelector.String()})
	ch := watcher.ResultChan()
	defer cancelFn()

	// iterate through all existing claims before streaming them
	claimList, err := iface.List(v1types.ListOptions{LabelSelector: claimLabelSelector.String()})
	if err != nil {
		return err
	}

	for _, wrapper := range claimList.Claims {
		evt := &Event{claim: wrapper, operation: watch.Added}
		receiveEvent(ctx, evt, iface, secretsNamespacer, lookup, lifecycler)
	}

	for {
		select {
		case evt, open := <-ch:
			// if the watch channel was closed, fail. this if statement is this is semantically
			// equivalent to if evt == nil {...}
			if !open {
				return errWatchClosed
			}
			logger.Debugf("received event %s", *evt.claim.Claim)
			receiveEvent(ctx, evt, iface, secretsNamespacer, lookup, lifecycler)
		case <-ctx.Done():
			logger.Debugf("loop has been stopped")
			return errLoopStopped
		}
	}
}

func receiveEvent(
	ctx context.Context,
	evt *Event,
	iface Interactor,
	secretsNamespacer v1.SecretsGetter,
	lookup k8s.ServiceCatalogLookup,
	lifecycler framework.Lifecycler,
) {
	nextAction, err := evt.nextAction()
	if isNoNextActionErr(err) {
		logger.Debugf("received event that has no next action (%s), skipping", err)
		return
	} else if err != nil {
		logger.Debugf("unknown error when determining the next action to make on the claim (%s)", err)
		return
	}

	claimUpdateCh := make(chan state.Update)
	wrapper := evt.claim

	cancelCtx, cancelFn := context.WithCancel(ctx)
	defer cancelFn()
	go nextAction(cancelCtx, evt, secretsNamespacer, lookup, lifecycler, claimUpdateCh)

	for {
		select {
		case claimUpdate := <-claimUpdateCh:
			logger.Debugf("wrapper before update: %s", *wrapper)
			state.UpdateClaim(wrapper.Claim, claimUpdate)
			logger.Debugf("wrapper after update: %s", *wrapper)

			w, err := iface.Update(wrapper)
			if err != nil {
				logger.Errorf("error updating claim %s (%s), stopping", wrapper.Claim.ClaimID, err)
				return
			}

			// if the claim update represents a terminal state, the loop should terminate
			if state.UpdateIsTerminal(claimUpdate) {
				return
			}
			wrapper = w
		case <-ctx.Done():
			return
		}
	}
}

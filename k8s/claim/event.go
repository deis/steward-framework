package claim

import (
	"context"
	"fmt"

	"github.com/deis/steward-framework"
	"github.com/deis/steward-framework/k8s"
	"github.com/deis/steward-framework/k8s/claim/state"
	"k8s.io/client-go/kubernetes/typed/core/v1"
	v1types "k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/pkg/watch"
)

type errNoNextAction struct {
	evt *Event
}

func (e errNoNextAction) Error() string {
	claim := e.evt.claim.Claim
	return fmt.Sprintf(
		"no next action for operation %s with event status %s, action %s",
		e.evt.operation,
		claim.Status,
		claim.Action,
	)
}

func isNoNextActionErr(e error) bool {
	_, ok := e.(errNoNextAction)
	return ok
}

type nextFunc func(
	context.Context,
	*Event,
	v1.SecretsGetter,
	k8s.ServiceCatalogLookup,
	framework.Lifecycler,
	chan<- state.Update,
)

// Event represents the event that a service plan claim has changed in kubernetes. It implements fmt.Stringer
type Event struct {
	claim     *k8s.ServicePlanClaimWrapper
	operation watch.EventType
}

func eventFromConfigMapEvent(raw watch.Event) (*Event, error) {
	configMap, ok := raw.Object.(*v1types.ConfigMap)
	if !ok {
		return nil, errNotAConfigMap
	}
	claimWrapper, err := k8s.ServicePlanClaimWrapperFromConfigMap(configMap)
	if err != nil {
		return nil, err
	}
	return &Event{
		claim:     claimWrapper,
		operation: raw.Type,
	}, nil
}

func (e Event) toConfigMap() *v1types.ConfigMap {
	return e.claim.ToConfigMap()
}

// String is the fmt.Stringer interface implementation
func (e Event) String() string {
	return fmt.Sprintf("%s %s", e.operation, *e.claim)
}

func (e *Event) nextAction() (nextFunc, error) {
	claim := e.claim.Claim
	action := k8s.ServicePlanClaimAction(claim.Action)
	status := k8s.ServicePlanClaimStatus(claim.Status)

	stateNoStatus := state.NewCurrentNoStatus(action, e.operation)
	stateWithStatus := state.NewCurrent(status, action, e.operation)

	nextFuncNoStatus, okNoStatus := transitionTable[stateNoStatus]
	nextFuncWithStatus, okWithStatus := transitionTable[stateWithStatus]
	if !okNoStatus && !okWithStatus {
		return nil, errNoNextAction{evt: e}
	} else if !okNoStatus {
		return nextFuncWithStatus, nil
	} else {
		return nextFuncNoStatus, nil
	}
}

package state

import (
	"fmt"

	"github.com/deis/steward-framework"
	"github.com/deis/steward-framework/k8s"
)

// ErrUpdate is an Update implementation that sets the claim to a failed state
type errUpdate struct {
	err error
}

// ErrUpdate returns a new Update implementation that has a failed status and status description equal to e.Error()
func ErrUpdate(e error) Update {
	return errUpdate{err: e}
}

func (e errUpdate) String() string {
	return fmt.Sprintf("status update to failure with error %s", e.err)
}

func (e errUpdate) Status() k8s.ServicePlanClaimStatus {
	return k8s.StatusFailed
}

func (e errUpdate) Description() string {
	return e.err.Error()
}

func (e errUpdate) InstanceID() string {
	return ""
}
func (e errUpdate) BindingID() string {
	return ""
}
func (e errUpdate) Extra() framework.JSONObject {
	return framework.EmptyJSONObject()
}

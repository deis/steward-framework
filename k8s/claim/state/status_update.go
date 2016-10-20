package state

import (
	"fmt"

	"github.com/deis/steward-framework"
	"github.com/deis/steward-framework/k8s"
)

type statusUpdate struct {
	status k8s.ServicePlanClaimStatus
}

// StatusUpdate returns an Update implementation that only updates the status field of a claim
func StatusUpdate(st k8s.ServicePlanClaimStatus) Update {
	return statusUpdate{status: st}
}

func (s statusUpdate) String() string {
	return fmt.Sprintf("status update to %s", s.Status)
}

func (s statusUpdate) Status() k8s.ServicePlanClaimStatus {
	return s.status
}

func (s statusUpdate) Description() string {
	return ""
}

func (s statusUpdate) InstanceID() string {
	return ""
}

func (s statusUpdate) BindingID() string {
	return ""
}

func (s statusUpdate) Extra() framework.JSONObject {
	return framework.EmptyJSONObject()
}

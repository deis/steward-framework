package state

import (
	"fmt"

	"github.com/deis/steward-framework"
	"github.com/deis/steward-framework/k8s"
)

type fullUpdate struct {
	status      k8s.ServicePlanClaimStatus
	description string
	instanceID  string
	bindingID   string
	extra       framework.JSONObject
}

// FullUpdate returns an Update implementation with all fields filled in
func FullUpdate(st k8s.ServicePlanClaimStatus, desc, instID, bindingID string, extra framework.JSONObject) Update {
	return fullUpdate{
		status:      st,
		description: desc,
		instanceID:  instID,
		bindingID:   bindingID,
		extra:       extra,
	}
}

func (f fullUpdate) String() string {
	return fmt.Sprintf(
		"full update. status = %s, description = '%s', instanceID = %s, bindingID = %s, extra = %s",
		f.status,
		f.description,
		f.instanceID,
		f.bindingID,
		f.extra,
	)
}

func (f fullUpdate) Status() k8s.ServicePlanClaimStatus {
	return f.status
}

func (f fullUpdate) Description() string {
	return f.description
}

func (f fullUpdate) InstanceID() string {
	return f.instanceID
}
func (f fullUpdate) BindingID() string {
	return f.bindingID
}
func (f fullUpdate) Extra() framework.JSONObject {
	return f.extra
}

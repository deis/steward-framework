package data

import (
	"strings"

	"k8s.io/client-go/pkg/api"
	"k8s.io/client-go/pkg/api/unversioned"
	"k8s.io/client-go/pkg/api/v1"
)

const (
	InstanceKind                           = "Instance"
	InstanceKindPlural                     = "Instances"
	InstanceStatePending     InstanceState = "Pending"
	InstanceStateProvisioned InstanceState = "Provisioned"
	InstanceStateFailed      InstanceState = "Failed"
)

// InstanceAPIResource returns an APIResource to describe the Instance third party resource
func InstanceAPIResource() *unversioned.APIResource {
	return &unversioned.APIResource{
		Name:       strings.ToLower(InstanceKindPlural),
		Namespaced: true,
		Kind:       InstanceKind,
	}
}

type Instance struct {
	unversioned.TypeMeta `json:",inline"`
	v1.ObjectMeta        `json:"metadata,omitempty"`

	Spec   InstanceSpec   `json:"spec"`
	Status InstanceStatus `json:"status"`
}

type InstanceSpec struct {
	ID              string              `json:"id"`
	ServiceClassRef api.ObjectReference `json:"service_class_ref"`
	// PlanID is the reference to the ServicePlan for this instance.
	PlanID string `json:"plan_id"`

	Parameters map[string]interface{} `json:"parameters"`
}

type InstanceStatus struct {
	Status InstanceState `json:"status"`
}

type InstanceState string

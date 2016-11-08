package data

import (
	"k8s.io/client-go/pkg/api"
	"k8s.io/client-go/pkg/api/unversioned"
	"k8s.io/client-go/pkg/api/v1"
)

type Instance struct {
	unversioned.TypeMeta `json:",inline"`
	v1.ObjectMeta        `json:"metadata,omitempty"`

	Spec   InstanceSpec   `json:"spec"`
	Status InstanceStatus `json:"status"`
}

type InstanceSpec struct {
	ID              string              `json:"id"`
	ServiceClassRef api.ObjectReference `json:"service_class_ref"`
	// PlanName is the reference to the ServicePlan for this instance.
	PlanName string

	Parameters map[string]interface{}
}

type InstanceStatus struct {
	Status InstanceState `json:"status"`
}

type InstanceState string

const (
	InstanceStatePending     InstanceState = "Pending"
	InstanceStateProvisioned InstanceState = "Provisioned"
	InstanceStateFailed      InstanceState = "Failed"
)

package data

import (
	"strings"

	"k8s.io/client-go/pkg/api"
	"k8s.io/client-go/pkg/api/unversioned"
	"k8s.io/client-go/pkg/api/v1"
)

const (
	ServiceInstanceKind                                  = "ServiceInstance"
	ServiceInstanceKindPlural                            = "ServiceInstances"
	ServiceInstanceStatePending     ServiceInstanceState = "Pending"
	ServiceInstanceStateProvisioned ServiceInstanceState = "Provisioned"
	ServiceInstanceStateFailed      ServiceInstanceState = "Failed"
)

// ServiceInstanceAPIResource returns an APIResource to describe the ServiceInstance third party resource
func ServiceInstanceAPIResource() *unversioned.APIResource {
	return &unversioned.APIResource{
		Name:       strings.ToLower(ServiceInstanceKindPlural),
		Namespaced: true,
		Kind:       ServiceInstanceKind,
	}
}

type ServiceInstance struct {
	unversioned.TypeMeta `json:",inline"`
	v1.ObjectMeta        `json:"metadata,omitempty"`

	Spec   ServiceInstanceSpec   `json:"spec"`
	Status ServiceInstanceStatus `json:"status"`
}

type ServiceInstanceSpec struct {
	ID              string              `json:"id"`
	ServiceClassRef api.ObjectReference `json:"service_class_ref"`
	// PlanID is the reference to the ServicePlan for this service instance.
	PlanID string `json:"plan_id"`

	Parameters map[string]interface{} `json:"parameters"`
}

type ServiceInstanceStatus struct {
	Status ServiceInstanceState `json:"status"`
}

type ServiceInstanceState string

package data

import (
	"strings"

	"github.com/deis/steward-framework/lib"
	"k8s.io/client-go/pkg/api"
	"k8s.io/client-go/pkg/api/unversioned"
)

const (
	ServiceBindingKind       = "ServiceBinding"
	ServiceBindingKindPlural = "ServiceBindings"
)

type ServiceBindingState string

const (
	ServiceBindingStatePending ServiceBindingState = "Pending"
	ServiceBindingStateBound   ServiceBindingState = "Bound"
	ServiceBindingStateFailed  ServiceBindingState = "Failed"
)

// ServiceBindingAPIResource returns an APIResource to describe the ServiceBinding third party
// resource
func ServiceBindingAPIResource() *unversioned.APIResource {
	return &unversioned.APIResource{
		Name:       strings.ToLower(ServiceBindingKindPlural),
		Namespaced: true,
		Kind:       ServiceBindingKind,
	}
}

type ServiceBinding struct {
	unversioned.TypeMeta `json:",inline"`
	api.ObjectMeta       `json:"metadata,omitempty"`

	Spec   ServiceBindingSpec   `json:"spec"`
	Status ServiceBindingStatus `json:"status"`
}

type ServiceBindingSpec struct {
	ID                 string              `json:"id"`
	ServiceInstanceRef api.ObjectReference `json:"service_instance_ref"`
	Parameters         lib.JSONObject      `json:"parameters"`
	SecretName         string              `json:"secret_name"`
}

type ServiceBindingStatus struct {
	State ServiceBindingState `json:"state"`
}

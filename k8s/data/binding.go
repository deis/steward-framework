package data

import (
	"github.com/deis/steward-framework/lib"
	"k8s.io/client-go/pkg/api"
	"k8s.io/client-go/pkg/api/unversioned"
)

const (
	BindingKind       = "Binding"
	BindingKindPlural = "Bindings"
)

type BindingState string

const (
	BindingStatePending BindingState = "Pending"
	BindingStateBound   BindingState = "Bound"
	BindingStateFailed  BindingState = "Failed"
)

type Binding struct {
	unversioned.TypeMeta
	api.ObjectMeta

	Spec   BindingSpec
	Status BindingStatus
}

type BindingSpec struct {
	ID          string              `json:"binding_id"`
	InstanceRef api.ObjectReference `json:"instance_ref"`
	Parameters  lib.JSONObject      `json:"parameters"`
	SecretName  string              `json:"secret_name"`
}

type BindingStatus struct {
	State BindingState
}

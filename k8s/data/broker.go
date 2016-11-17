package data

import (
	"strings"

	"github.com/deis/steward-framework"
	"k8s.io/client-go/pkg/api/unversioned"
	"k8s.io/client-go/pkg/api/v1"
)

const (
	BrokerKind       = "Broker"
	BrokerKindPlural = "Brokers"
)

// BrokerAPIResource returns an APIResource to describe the Broker third party resource
func BrokerAPIResource() *unversioned.APIResource {
	return &unversioned.APIResource{
		Name:       strings.ToLower(BrokerKindPlural),
		Namespaced: true,
		Kind:       BrokerKind,
	}
}

type BrokerState string

const (
	BrokerStatePending   BrokerState = "Pending"
	BrokerStateAvailable BrokerState = "Available"
	BrokerStateFailed    BrokerState = "Failed"
)

type Broker struct {
	unversioned.TypeMeta `json:",inline"`
	v1.ObjectMeta        `json:"metadata,omitempty"`

	Spec   framework.BrokerSpec
	Status BrokerStatus
}

type BrokerStatus struct {
	State BrokerState `json:"state"`
}

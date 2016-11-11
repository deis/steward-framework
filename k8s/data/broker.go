package data

import (
	"github.com/deis/steward-framework"
	"k8s.io/client-go/pkg/api/unversioned"
	"k8s.io/client-go/pkg/api/v1"
)

const (
	BrokerKind       = "Broker"
	BrokerKindPlural = "Brokers"
)

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

package data

import (
	"k8s.io/client-go/pkg/api/unversioned"
	"k8s.io/client-go/pkg/api/v1"
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

	Spec   BrokerSpec
	Status BrokerStatus
}

type BrokerSpec struct {
	URL      string `json:"url"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type BrokerStatus struct {
	State BrokerState `json:"state"`
}

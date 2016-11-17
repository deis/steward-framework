package broker

import (
	"github.com/deis/steward-framework/k8s/data"
	"k8s.io/client-go/dynamic"
)

// UpdateBrokerFunc is the function that can update a broker
type UpdateBrokerFunc func(*data.Broker) (*data.Broker, error)

// NewK8sUpdateBrokerFunc returns an UpdateBrokerFunc backed by a Kubernetes client
func NewK8sUpdateBrokerFunc(cl *dynamic.Client) UpdateBrokerFunc {
	return func(newBroker *data.Broker) (*data.Broker, error) {
		resCl := cl.Resource(data.BrokerAPIResource(), newBroker.Namespace)
		unstruc, err := data.TranslateToUnstructured(newBroker)
		if err != nil {
			return nil, err
		}
		retUnstruc, err := resCl.Update(unstruc)
		if err != nil {
			return nil, err
		}
		retBroker := new(data.Broker)
		if err := data.TranslateToTPR(retUnstruc, retBroker, data.BrokerKind); err != nil {
			return nil, err
		}
		return retBroker, nil
	}
}

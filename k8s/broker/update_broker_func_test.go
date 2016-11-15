package broker

import (
	"github.com/deis/steward-framework/k8s/data"
)

func newFakeUpdateBrokerFunc(retErr error) (UpdateBrokerFunc, *[]*data.Broker) {
	brokers := new([]*data.Broker)
	return func(newBroker *data.Broker) (*data.Broker, error) {
		if retErr != nil {
			return nil, retErr
		}
		*brokers = append(*brokers, newBroker)
		return newBroker, nil
	}, brokers
}

package servicebroker

import (
	"github.com/deis/steward-framework/k8s/data"
)

func newFakeUpdateServiceBrokerFunc(retErr error) (UpdateServiceBrokerFunc, *[]*data.ServiceBroker) {
	serviceBrokers := new([]*data.ServiceBroker)
	return func(newServiceBroker *data.ServiceBroker) (*data.ServiceBroker, error) {
		if retErr != nil {
			return nil, retErr
		}
		*serviceBrokers = append(*serviceBrokers, newServiceBroker)
		return newServiceBroker, nil
	}, serviceBrokers
}

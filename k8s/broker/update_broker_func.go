package broker

import (
	"github.com/deis/steward-framework/k8s/data"
	"github.com/deis/steward-framework/k8s/restutil"
	"k8s.io/client-go/rest"
)

// UpdateBrokerFunc is the function that can update a broker
type UpdateBrokerFunc func(*data.Broker) (*data.Broker, error)

// NewK8sUpdateBrokerFunc returns an UpdateBrokerFunc backed by a Kubernetes client
func NewK8sUpdateBrokerFunc(restIface rest.Interface) UpdateBrokerFunc {
	return func(newBroker *data.Broker) (*data.Broker, error) {
		url := restutil.AbsPath(
			restutil.APIVersionBase,
			restutil.APIVersion,
			true,
			newBroker.Namespace,
			data.BrokerKindPlural,
		)
		res := &data.Broker{}
		if err := restIface.Put().AbsPath(url...).Body(newBroker).Do().Into(res); err != nil {
			return nil, err
		}
		return res, nil
	}
}

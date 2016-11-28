package refs

import (
	"github.com/deis/steward-framework/k8s/data"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/pkg/api"
)

// BrokerGetterFunc is the function that attempts to retrieve a broker at the given object ref
type BrokerGetterFunc func(api.ObjectReference) (*data.Broker, error)

// NewK8sBrokerGetterFunc returns a BrokerGetterFunc  backed by a real kubernetes client
func NewK8sBrokerGetterFunc(cl *dynamic.Client) BrokerGetterFunc {
	return func(ref api.ObjectReference) (*data.Broker, error) {
		resCl := cl.Resource(data.BrokerAPIResource(), ref.Namespace)
		unstruc, err := resCl.Get(ref.Name)
		if err != nil {
			return nil, err
		}
		retBroker := new(data.Broker)
		if err := data.TranslateToTPR(unstruc, retBroker, data.BrokerKind); err != nil {
			return nil, err
		}
		return retBroker, nil
	}
}

// NewFakeBrokerGetterFunc returns a fake BrokerGetterFunc. If retErr is non-nil, it always returns
// (nil, retErr). Otherwise returns (broker, nil)
func NewFakeBrokerGetterFunc(broker *data.Broker, retErr error) BrokerGetterFunc {
	return func(api.ObjectReference) (*data.Broker, error) {
		if retErr != nil {
			return nil, retErr
		}
		return broker, nil
	}
}

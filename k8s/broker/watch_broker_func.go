package broker

import (
	"github.com/deis/steward-framework/k8s/data"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/pkg/watch"
)

// WatchBrokerFunc is the function that returns a watch interface for broker resources
type WatchBrokerFunc func(namespace string) (watch.Interface, error)

// NewK8sWatchBrokerFunc returns a WatchBrokerFunc backed by a Kubernetes client
func NewK8sWatchBrokerFunc(cl *dynamic.Client) WatchBrokerFunc {
	return func(namespace string) (watch.Interface, error) {
		resCl := cl.Resource(data.BrokerAPIResource(), namespace)
		return resCl.Watch(&data.Broker{})
	}
}

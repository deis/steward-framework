package broker

import (
	"strings"

	"github.com/deis/steward-framework/k8s/data"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/pkg/api/unversioned"
	"k8s.io/client-go/pkg/watch"
)

// WatchBrokerFunc is the function that returns a watch interface for broker resources
type WatchBrokerFunc func(namespace string) (watch.Interface, error)

var brokerAPIResource = unversioned.APIResource{
	Name:       strings.ToLower(data.BrokerKindPlural),
	Namespaced: true,
	Kind:       data.BrokerKind,
}

// NewK8sWatchBrokerFunc returns a WatchBrokerFunc backed by a Kubernetes client
func NewK8sWatchBrokerFunc(cl *dynamic.Client) WatchBrokerFunc {
	return func(namespace string) (watch.Interface, error) {
		resCl := cl.Resource(&brokerAPIResource, namespace)
		return resCl.Watch(&data.Broker{})
	}
}

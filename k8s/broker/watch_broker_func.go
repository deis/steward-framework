package broker

import (
	"github.com/deis/steward-framework/k8s/data"
	"github.com/deis/steward-framework/k8s/restutil"
	"k8s.io/client-go/pkg/watch"
	"k8s.io/client-go/rest"
)

// WatchBrokerFunc is the function that returns a watch interface for broker resources
type WatchBrokerFunc func(namespace string) (watch.Interface, error)

// NewK8sWatchBrokerFunc returns a WatchBrokerFunc backed by a Kubernetes client
func NewK8sWatchBrokerFunc(restIface rest.Interface) WatchBrokerFunc {
	return func(namespace string) (watch.Interface, error) {
		url := restutil.AbsPath(
			restutil.APIVersionBase,
			restutil.APIVersion,
			true,
			namespace,
			data.BrokerKindPlural,
		)
		return restIface.Get().AbsPath(url...).Watch()
	}
}

package instance

import (
	"github.com/deis/steward-framework/k8s/data"
	"github.com/deis/steward-framework/k8s/restutil"
	"k8s.io/client-go/pkg/watch"
	"k8s.io/client-go/rest"
)

// WatchInstanceFunc is the function that returns a watch interface for instance resources
type WatchInstanceFunc func(namespace string) (watch.Interface, error)

// NewK8sWatchInstanceFunc returns a WatchInstanceFunc backed by a Kubernetes client
func NewK8sWatchInstanceFunc(restIface rest.Interface) WatchInstanceFunc {
	return func(namespace string) (watch.Interface, error) {
		url := restutil.AbsPath(
			restutil.APIVersionBase,
			restutil.APIVersion,
			true,
			namespace,
			data.InstanceKindPlural,
		)
		return restIface.Get().AbsPath(url...).Watch()
	}
}

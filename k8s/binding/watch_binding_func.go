package binding

import (
	"github.com/deis/steward-framework/k8s/data"
	"github.com/deis/steward-framework/k8s/restutil"
	"k8s.io/client-go/pkg/watch"
	"k8s.io/client-go/rest"
)

// WatchBindingFunc is the function that returns a watch interface for binding resources
type WatchBindingFunc func(namespace string) (watch.Interface, error)

// NewK8sWatchBindingFunc returns a WatchBindingFunc backed by a Kubernetes client
func NewK8sWatchBindingFunc(restIface rest.Interface) WatchBindingFunc {
	return func(namespace string) (watch.Interface, error) {
		url := restutil.AbsPath(
			restutil.APIVersionBase,
			restutil.APIVersion,
			true,
			namespace,
			data.BindingKindPlural,
		)
		return restIface.Get().AbsPath(url...).Watch()
	}
}

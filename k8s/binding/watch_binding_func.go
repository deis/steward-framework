package binding

import (
	"github.com/deis/steward-framework/k8s/data"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/pkg/api/unversioned"
	"k8s.io/client-go/pkg/watch"
)

// WatchBindingFunc is the function that returns a watch interface for binding resources
type WatchBindingFunc func(namespace string) (watch.Interface, error)

var bindingAPIResource = unversioned.APIResource{
	Name:       "",
	Namespaced: true,
	Kind:       data.BindingKindPlural,
}

// NewK8sWatchBindingFunc returns a WatchBindingFunc backed by a Kubernetes client
func NewK8sWatchBindingFunc(cl *dynamic.Client) WatchBindingFunc {
	return func(namespace string) (watch.Interface, error) {
		resCl := cl.Resource(&bindingAPIResource, namespace)
		return resCl.Watch(&data.Binding{})
	}
}

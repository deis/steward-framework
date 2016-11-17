package binding

import (
	"github.com/deis/steward-framework/k8s/data"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/pkg/watch"
)

// WatchBindingFunc is the function that returns a watch interface for binding resources
type WatchBindingFunc func(namespace string) (watch.Interface, error)

// NewK8sWatchBindingFunc returns a WatchBindingFunc backed by a Kubernetes client
func NewK8sWatchBindingFunc(cl *dynamic.Client) WatchBindingFunc {
	return func(namespace string) (watch.Interface, error) {
		resCl := cl.Resource(data.BindingAPIResource(), namespace)
		// TODO: call watch.Filter here, and call data.TranslateToTPR in the filter func.
		// Do this so the loop doesn't have to call it and instead can just type-assert
		return resCl.Watch(&data.Binding{})
	}
}

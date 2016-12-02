package servicebinding

import (
	"github.com/deis/steward-framework/k8s/data"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/pkg/watch"
)

// WatchServiceBindingFunc is the function that returns a watch interface for service binding
// resources
type WatchServiceBindingFunc func(namespace string) (watch.Interface, error)

// NewK8sWatchServiceBindingFunc returns a WatchServiceBindingFunc backed by a Kubernetes client
func NewK8sWatchServiceBindingFunc(cl *dynamic.Client) WatchServiceBindingFunc {
	return func(namespace string) (watch.Interface, error) {
		resCl := cl.Resource(data.ServiceBindingAPIResource(), namespace)
		// TODO: call watch.Filter here, and call data.TranslateToTPR in the filter func.
		// Do this so the loop doesn't have to call it and instead can just type-assert
		return resCl.Watch(&data.ServiceBinding{})
	}
}

package instance

import (
	"github.com/deis/steward-framework/k8s/data"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/pkg/watch"
)

// WatchInstanceFunc is the function that returns a watch interface for instance resources
type WatchInstanceFunc func(namespace string) (watch.Interface, error)

// NewK8sWatchInstanceFunc returns a WatchInstanceFunc backed by a Kubernetes client
func NewK8sWatchInstanceFunc(cl *dynamic.Client) WatchInstanceFunc {
	return func(namespace string) (watch.Interface, error) {
		resCl := cl.Resource(data.InstanceAPIResource(), namespace)
		// TODO: call watch.Filter here, and call data.TranslateToTPR in the filter func.
		// Do this so the loop doesn't have to call it and instead can just type-assert
		return resCl.Watch(&data.Instance{})
	}
}

package serviceinstance

import (
	"github.com/deis/steward-framework/k8s/data"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/pkg/watch"
)

// WatchServiceInstanceFunc is the function that returns a watch interface for service instance
// resources
type WatchServiceInstanceFunc func(namespace string) (watch.Interface, error)

// NewK8sWatchServiceInstanceFunc returns a WatchServiceInstanceFunc backed by a Kubernetes client
func NewK8sWatchServiceInstanceFunc(cl *dynamic.Client) WatchServiceInstanceFunc {
	return func(namespace string) (watch.Interface, error) {
		resCl := cl.Resource(data.ServiceInstanceAPIResource(), namespace)
		// TODO: call watch.Filter here, and call data.TranslateToTPR in the filter func.
		// Do this so the loop doesn't have to call it and instead can just type-assert
		return resCl.Watch(&data.ServiceInstance{})
	}
}

package instance

import (
	"strings"

	"github.com/deis/steward-framework/k8s/data"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/pkg/api/unversioned"
	"k8s.io/client-go/pkg/watch"
)

// WatchInstanceFunc is the function that returns a watch interface for instance resources
type WatchInstanceFunc func(namespace string) (watch.Interface, error)

var instanceAPIResource = unversioned.APIResource{
	Name:       strings.ToLower(data.InstanceKindPlural),
	Namespaced: true,
	Kind:       data.InstanceKind,
}

// NewK8sWatchInstanceFunc returns a WatchInstanceFunc backed by a Kubernetes client
func NewK8sWatchInstanceFunc(cl *dynamic.Client) WatchInstanceFunc {
	return func(namespace string) (watch.Interface, error) {
		resCl := cl.Resource(&instanceAPIResource, namespace)
		return resCl.Watch(&data.Instance{})
	}
}

package refs

import (
	"github.com/deis/steward-framework/k8s/data"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/pkg/api"
)

// ServiceInstanceGetterFunc is the function that attempts to fetch a service instance at the given
// object ref
type ServiceInstanceGetterFunc func(api.ObjectReference) (*data.ServiceInstance, error)

// NewK8sServiceInstanceGetterFunc returns an ServiceInstanceGetterFunc backed by a real kubernetes
// client
func NewK8sServiceInstanceGetterFunc(cl *dynamic.Client) ServiceInstanceGetterFunc {
	return func(ref api.ObjectReference) (*data.ServiceInstance, error) {
		resCl := cl.Resource(data.ServiceInstanceAPIResource(), ref.Namespace)
		unstruc, err := resCl.Get(ref.Name)
		if err != nil {
			return nil, err
		}
		retServiceInstance := new(data.ServiceInstance)
		if err := data.TranslateToTPR(unstruc, retServiceInstance, data.ServiceInstanceKind); err != nil {
			return nil, err
		}
		return retServiceInstance, nil
	}
}

// NewFakeServiceInstanceGetterFunc returns a fake ServiceInstanceGetterFunc. If retErr is non-nil,
// it always returns (nil, retErr). Otherwise returns (inst, nil)
func NewFakeServiceInstanceGetterFunc(inst *data.ServiceInstance, retErr error) ServiceInstanceGetterFunc {
	return func(api.ObjectReference) (*data.ServiceInstance, error) {
		if retErr != nil {
			return nil, retErr
		}
		return inst, nil
	}
}

package serviceinstance

import (
	"github.com/deis/steward-framework/k8s/data"
	"k8s.io/client-go/dynamic"
)

// UpdateServiceInstanceFunc is the function that can update a service instance
type UpdateServiceInstanceFunc func(*data.ServiceInstance) (*data.ServiceInstance, error)

// NewK8sUpdateServiceInstanceFunc returns an UpdateServiceInstanceFunc backed by a Kubernetes
// client
func NewK8sUpdateServiceInstanceFunc(cl *dynamic.Client) UpdateServiceInstanceFunc {
	return func(newServiceInstance *data.ServiceInstance) (*data.ServiceInstance, error) {
		resCl := cl.Resource(data.ServiceInstanceAPIResource(), newServiceInstance.Namespace)
		unstruc, err := data.TranslateToUnstructured(newServiceInstance)
		if err != nil {
			return nil, err
		}
		retUnstruc, err := resCl.Update(unstruc)
		if err != nil {
			return nil, err
		}
		retServiceInstance := new(data.ServiceInstance)
		if err := data.TranslateToTPR(retUnstruc, retServiceInstance, data.ServiceInstanceKind); err != nil {
			return nil, err
		}
		return retServiceInstance, nil
	}
}

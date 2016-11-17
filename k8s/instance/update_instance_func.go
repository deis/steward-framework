package instance

import (
	"github.com/deis/steward-framework/k8s/data"
	"k8s.io/client-go/dynamic"
)

// UpdateInstanceFunc is the function that can update an instance
type UpdateInstanceFunc func(*data.Instance) (*data.Instance, error)

// NewK8sUpdateInstanceFunc returns an UpdateInstanceFunc backed by a Kubernetes client
func NewK8sUpdateInstanceFunc(cl *dynamic.Client) UpdateInstanceFunc {
	return func(newInstance *data.Instance) (*data.Instance, error) {
		resCl := cl.Resource(&instanceAPIResource, newInstance.Namespace)
		unstruc, err := data.TranslateToUnstructured(newInstance)
		if err != nil {
			return nil, err
		}
		retUnstruc, err := resCl.Update(unstruc)
		if err != nil {
			return nil, err
		}
		retInstance := new(data.Instance)
		if err := data.TranslateToTPR(retUnstruc, retInstance, data.InstanceKind); err != nil {
			return nil, err
		}
		return retInstance, nil
	}
}

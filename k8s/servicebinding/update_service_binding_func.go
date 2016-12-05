package servicebinding

import (
	"github.com/deis/steward-framework/k8s/data"
	"k8s.io/client-go/dynamic"
)

// UpdateServiceBindingFunc is the function that can update a service binding
type UpdateServiceBindingFunc func(*data.ServiceBinding) (*data.ServiceBinding, error)

// NewK8sUpdateServiceBindingFunc returns an UpdateServiceBindingFunc backed by a Kubernetes client
func NewK8sUpdateServiceBindingFunc(cl *dynamic.Client) UpdateServiceBindingFunc {
	return func(newServiceBinding *data.ServiceBinding) (*data.ServiceBinding, error) {
		resCl := cl.Resource(data.ServiceBindingAPIResource(), newServiceBinding.Namespace)
		unstruc, err := data.TranslateToUnstructured(newServiceBinding)
		if err != nil {
			return nil, err
		}
		retUnstruc, err := resCl.Update(unstruc)
		if err != nil {
			return nil, err
		}
		retServiceBinding := new(data.ServiceBinding)
		if err := data.TranslateToTPR(retUnstruc, retServiceBinding, data.ServiceBindingKind); err != nil {
			return nil, err
		}
		return retServiceBinding, nil
	}
}

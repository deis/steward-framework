package binding

import (
	"github.com/deis/steward-framework/k8s/data"
	"github.com/deis/steward-framework/k8s/restutil"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
)

// UpdateBindingFunc is the function that can update an instance
type UpdateBindingFunc func(*data.Binding) (*data.Binding, error)

// NewK8sUpdateInstanceFunc returns an UpdateInstanceFunc backed by a Kubernetes client
func NewK8sUpdateBindingFunc(cl *dynamic.Client) UpdateBindingFunc {
	return func(newBinding *data.Binding) (*data.Binding, error) {
		resCl := resCl := cl.Resource(&bindingAPIResource, newBinding.Namespace)
		unstruc, err := data.ToUnstructured(newBinding)
		if err != nil {
			return nil, err
		}
		retUnstruc, err := resCl.Update(unstruc)
		if err != nil {
			return nil, err
		}
		retBinding := new(data.Binding)
		if err := data.TranslateToTPR(retUnstruc, retBinding); err != nil {
			return nil, err
		}
		return retBinding, nil
	}
}

package binding

import (
	"github.com/deis/steward-framework/k8s/data"
	"k8s.io/client-go/dynamic"
)

// UpdateBindingFunc is the function that can update a binding
type UpdateBindingFunc func(*data.Binding) (*data.Binding, error)

// NewK8sUpdateBindingFunc returns an UpdateBindingFunc backed by a Kubernetes client
func NewK8sUpdateBindingFunc(cl *dynamic.Client) UpdateBindingFunc {
	return func(newBinding *data.Binding) (*data.Binding, error) {
		resCl := cl.Resource(data.BindingAPIResource(), newBinding.Namespace)
		unstruc, err := data.TranslateToUnstructured(newBinding)
		if err != nil {
			return nil, err
		}
		retUnstruc, err := resCl.Update(unstruc)
		if err != nil {
			return nil, err
		}
		retBinding := new(data.Binding)
		if err := data.TranslateToTPR(retUnstruc, retBinding, data.BindingKind); err != nil {
			return nil, err
		}
		return retBinding, nil
	}
}

package binding

import (
	"github.com/deis/steward-framework/k8s/data"
)

func newFakeUpdateBindingFunc(retErr error) (UpdateBindingFunc, *[]*data.Binding) {
	bindings := new([]*data.Binding)
	return func(newBinding *data.Binding) (*data.Binding, error) {
		if retErr != nil {
			return nil, retErr
		}
		*bindings = append(*bindings, newBinding)
		return newBinding, nil
	}, bindings
}

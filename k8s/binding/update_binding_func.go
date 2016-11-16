package binding

import (
	"github.com/deis/steward-framework/k8s/data"
	"github.com/deis/steward-framework/k8s/restutil"
	"k8s.io/client-go/rest"
)

// UpdateBindingFunc is the function that can update an instance
type UpdateBindingFunc func(*data.Binding) (*data.Binding, error)

// NewK8sUpdateInstanceFunc returns an UpdateInstanceFunc backed by a Kubernetes client
func NewK8sUpdateBindingFunc(restIface rest.Interface) UpdateBindingFunc {
	return func(newBinding *data.Binding) (*data.Binding, error) {
		url := restutil.AbsPath(
			restutil.APIVersionBase,
			restutil.APIVersion,
			true,
			newBinding.Namespace,
			data.BindingKindPlural,
		)
		res := &data.Binding{}
		if err := restIface.Put().AbsPath(url...).Body(newBinding).Do().Into(res); err != nil {
			return nil, err
		}
		return res, nil
	}
}

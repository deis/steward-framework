package servicebinding

import (
	"github.com/deis/steward-framework/k8s/data"
)

func newFakeUpdateServiceBindingFunc(retErr error) (UpdateServiceBindingFunc, *[]*data.ServiceBinding) {
	serviceBindings := new([]*data.ServiceBinding)
	return func(newServiceBinding *data.ServiceBinding) (*data.ServiceBinding, error) {
		if retErr != nil {
			return nil, retErr
		}
		*serviceBindings = append(*serviceBindings, newServiceBinding)
		return newServiceBinding, nil
	}, serviceBindings
}

package broker

import (
	"github.com/deis/steward-framework/k8s/data"
	"github.com/deis/steward-framework/k8s/restutil"
	"k8s.io/client-go/rest"
)

// CreateServiceClassFunc is the function that can successfully create a ServiceClass
type CreateServiceClassFunc func(*data.ServiceClass) error

// NewK8sCreateServiceClassFunc returns a CreateServiceClassFunc implemented with restIFace
func NewK8sCreateServiceClassFunc(restIface rest.Interface) CreateServiceClassFunc {
	return func(sClass *data.ServiceClass) error {
		url := restutil.AbsPath(
			restutil.APIVersionBase,
			restutil.APIVersion,
			false,
			sClass.Namespace,
			data.ServiceClassKindPlural,
		)
		return restIface.Post().AbsPath(url...).Do().Error()
	}
}

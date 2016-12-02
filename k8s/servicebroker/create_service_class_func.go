package servicebroker

import (
	"github.com/deis/steward-framework/k8s/data"
	"k8s.io/client-go/dynamic"
)

// CreateServiceClassFunc is the function that can successfully create a ServiceClass
type CreateServiceClassFunc func(*data.ServiceClass) error

// NewK8sCreateServiceClassFunc returns a CreateServiceClassFunc implemented with restIFace
func NewK8sCreateServiceClassFunc(cl *dynamic.Client) CreateServiceClassFunc {
	return func(sClass *data.ServiceClass) error {
		resCl := cl.Resource(data.ServiceClassAPIResource(), sClass.Namespace)
		unstruc, err := data.TranslateToUnstructured(sClass)
		if err != nil {
			return err
		}
		retUnstruc, err := resCl.Create(unstruc)
		if err != nil {
			return err
		}
		retSClass := new(data.ServiceClass)
		if err := data.TranslateToTPR(retUnstruc, retSClass, data.ServiceClassKind); err != nil {
			return err
		}
		return nil
	}
}

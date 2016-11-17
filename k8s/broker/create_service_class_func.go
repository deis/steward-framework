package broker

import (
	"strings"

	"github.com/deis/steward-framework/k8s/data"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/pkg/api/unversioned"
)

// CreateServiceClassFunc is the function that can successfully create a ServiceClass
type CreateServiceClassFunc func(*data.ServiceClass) error

var serviceClassAPIResource = unversioned.APIResource{
	Name:       strings.ToLower(data.ServiceClassKindPlural),
	Namespaced: true,
	Kind:       data.ServiceClassKind,
}

// NewK8sCreateServiceClassFunc returns a CreateServiceClassFunc implemented with restIFace
func NewK8sCreateServiceClassFunc(cl *dynamic.Client) CreateServiceClassFunc {
	return func(sClass *data.ServiceClass) error {
		resCl := cl.Resource(&serviceClassAPIResource, sClass.Namespace)
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

package refs

import (
	"github.com/deis/steward-framework/k8s/data"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/pkg/api"
)

// ServiceClassGetterFunc is the function that attempts to retrieve a service class at the
// given object ref
type ServiceClassGetterFunc func(api.ObjectReference) (*data.ServiceClass, error)

// NewK8sServiceClassGetterFunc returns a ServiceClassGetterFunc backed by a real kubernetes client
func NewK8sServiceClassGetterFunc(cl *dynamic.Client) ServiceClassGetterFunc {
	return func(ref api.ObjectReference) (*data.ServiceClass, error) {
		resCl := cl.Resource(data.ServiceClassAPIResource(), ref.Namespace)
		unstruc, err := resCl.Get(ref.Name)
		if err != nil {
			return nil, err
		}
		retSvcClass := new(data.ServiceClass)
		if err := data.TranslateToTPR(unstruc, retSvcClass, data.ServiceClassKind); err != nil {
			return nil, err
		}
		return retSvcClass, nil
	}
}

// NewFakeServiceClassGetterFunc returns a fake ServiceClassGetterFunc. If retErr is non-nil,
// it always returns (nil, retErr). Otherwise returns (svcClass, nil)
func NewFakeServiceClassGetterFunc(svcClass *data.ServiceClass, retErr error) ServiceClassGetterFunc {
	return func(api.ObjectReference) (*data.ServiceClass, error) {
		if retErr != nil {
			return nil, retErr
		}
		return svcClass, nil
	}
}

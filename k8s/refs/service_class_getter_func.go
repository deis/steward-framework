package refs

import (
	"github.com/deis/steward-framework/k8s/data"
	"github.com/deis/steward-framework/k8s/restutil"
	"k8s.io/client-go/pkg/api"
	"k8s.io/client-go/rest"
)

// ServiceClassGetterFunc is the function that attempts to retrieve a service class at the
// given object ref
type ServiceClassGetterFunc func(api.ObjectReference) (*data.ServiceClass, error)

// NewK8sServiceClassGetterFunc returns a ServiceClassGetterFunc backed by a real kubernetes client
func NewK8sServiceClassGetterFunc(restIface rest.Interface) ServiceClassGetterFunc {
	return func(ref api.ObjectReference) (*data.ServiceClass, error) {
		ret := new(data.ServiceClass)
		url := append(
			restutil.AbsPath(
				restutil.APIVersionBase,
				restutil.APIVersion,
				false,
				ref.Namespace,
				data.ServiceClassKindPlural,
			),
			ref.Name,
		)
		if err := restIface.Get().AbsPath(url...).Do().Into(ret); err != nil {
			return nil, err
		}
		return ret, nil
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

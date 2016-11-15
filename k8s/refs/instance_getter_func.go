package refs

import (
	"github.com/deis/steward-framework/k8s/data"
	"github.com/deis/steward-framework/k8s/restutil"
	"k8s.io/client-go/pkg/api"
	"k8s.io/client-go/rest"
)

// InstanceGetterFunc is the function that attempts to fetch an instance at the given object ref
type InstanceGetterFunc func(api.ObjectReference) (*data.Instance, error)

// NewK8sInstanceGetterFunc returns an InstanceGetterFunc backed by a real kubernetes client
func NewK8sInstanceGetterFunc(restIface rest.Interface) InstanceGetterFunc {
	return func(ref api.ObjectReference) (*data.Instance, error) {
		ret := new(data.Instance)
		url := append(
			restutil.AbsPath(
				restutil.APIVersionBase,
				restutil.APIVersion,
				false,
				ref.Namespace,
				data.InstanceKindPlural,
			),
			ref.Name,
		)

		if err := restIface.Get().AbsPath(url...).Do().Into(ret); err != nil {
			return nil, err
		}
		return ret, nil
	}
}

// NewFakeInstanceGetterFunc returns a fake InstanceGetterFunc. If retErr is non-nil, it always returns
// (nil, retErr). Otherwise returns (inst, nil)
func NewFakeInstanceGetterFunc(inst *data.Instance, retErr error) InstanceGetterFunc {
	return func(api.ObjectReference) (*data.Instance, error) {
		if retErr != nil {
			return nil, retErr
		}
		return inst, nil
	}
}

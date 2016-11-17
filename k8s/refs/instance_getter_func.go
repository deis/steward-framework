package refs

import (
	"github.com/deis/steward-framework/k8s/data"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/pkg/api"
)

// InstanceGetterFunc is the function that attempts to fetch an instance at the given object ref
type InstanceGetterFunc func(api.ObjectReference) (*data.Instance, error)

// NewK8sInstanceGetterFunc returns an InstanceGetterFunc backed by a real kubernetes client
func NewK8sInstanceGetterFunc(cl *dynamic.Client) InstanceGetterFunc {
	return func(ref api.ObjectReference) (*data.Instance, error) {
		resCl := cl.Resource(data.InstanceAPIResource(), ref.Namespace)
		unstruc, err := resCl.Get(ref.Name)
		if err != nil {
			return nil, err
		}
		retInstance := new(data.Instance)
		if err := data.TranslateToTPR(unstruc, retInstance, data.InstanceKind); err != nil {
			return nil, err
		}
		return retInstance, nil
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

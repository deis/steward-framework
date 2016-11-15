package instance

import (
	"github.com/deis/steward-framework/k8s/data"
	"github.com/deis/steward-framework/k8s/restutil"
	"k8s.io/client-go/rest"
)

// UpdateInstancdFunc is the function that can update an instance
type UpdateInstanceFunc func(*data.Instance) (*data.Instance, error)

// NewK8sUpdateInstanceFunc returns an UpdateInstanceFunc backed by a Kubernetes client
func NewK8sUpdateInstanceFunc(restIface rest.Interface) UpdateInstanceFunc {
	return func(newInstance *data.Instance) (*data.Instance, error) {
		url := restutil.AbsPath(
			restutil.APIVersionBase,
			restutil.APIVersion,
			true,
			newInstance.Namespace,
			data.InstanceKindPlural,
		)
		res := &data.Instance{}
		if err := restIface.Put().AbsPath(url...).Body(newInstance).Do().Into(res); err != nil {
			return nil, err
		}
		return res, nil
	}
}

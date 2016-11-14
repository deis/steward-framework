package instance

import (
	"github.com/deis/steward-framework/k8s/data"
)

func newFakeUpdateInstanceFunc(retErr error) (UpdateInstanceFunc, *[]*data.Instance) {
	instances := new([]*data.Instance)
	return func(newInstance *data.Instance) (*data.Instance, error) {
		if retErr != nil {
			return nil, retErr
		}
		*instances = append(*instances, newInstance)
		return newInstance, nil
	}, instances
}

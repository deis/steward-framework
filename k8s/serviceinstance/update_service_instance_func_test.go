package serviceinstance

import (
	"github.com/deis/steward-framework/k8s/data"
)

func newFakeUpdateServiceInstanceFunc(retErr error) (UpdateServiceInstanceFunc, *[]*data.ServiceInstance) {
	serviceInstances := new([]*data.ServiceInstance)
	return func(newServiceInstance *data.ServiceInstance) (*data.ServiceInstance, error) {
		if retErr != nil {
			return nil, retErr
		}
		*serviceInstances = append(*serviceInstances, newServiceInstance)
		return newServiceInstance, nil
	}, serviceInstances
}

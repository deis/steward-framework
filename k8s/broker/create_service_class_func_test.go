package broker

import (
	"github.com/deis/steward-framework/k8s/data"
)

// returns the function and a mutable slice of classes that were created. if retErr != nil,
// it is always returned by the function and the returned slice is never modified
func newFakeCreateServiceClassFunc(retErr error) (CreateServiceClassFunc, *[]*data.ServiceClass) {
	createdClasses := []*data.ServiceClass{}
	retFn := func(sClass *data.ServiceClass) error {
		if retErr != nil {
			return retErr
		}
		createdClasses = append(createdClasses, sClass)
		return nil
	}
	return retFn, &createdClasses
}

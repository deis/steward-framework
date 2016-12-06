package fake

import (
	"context"

	"github.com/deis/steward-framework"
)

// OperationStatusRetriever is a fake implementation of the
// (github.com/deis/steward-framework).OperationStatusRetriever interface that is suitable for use
// in unit tests.
type OperationStatusRetriever struct {
	Reqs []*framework.OperationStatusRequest
	Res  func() *framework.OperationStatusResponse
}

// GetOperationStatus is the OperationStatusRetriever implementation. It just records the
// OperationStatusRequest and returns a hardcoded OperationStatusResponse.
func (o *OperationStatusRetriever) GetOperationStatus(
	ctx context.Context,
	brokerSpec framework.ServiceBrokerSpec,
	req *framework.OperationStatusRequest,
) (*framework.OperationStatusResponse, error) {

	o.Reqs = append(o.Reqs, req)
	return o.Res(), nil
}

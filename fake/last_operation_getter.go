package fake

import (
	"context"

	"github.com/deis/steward-framework"
)

// LastOperationGetter is a fake implementation of the
// (github.com/deis/steward-framework).LastOperationGetter interface that is suitable for use in
// unit tests.
type LastOperationGetter struct {
	Reqs []*framework.GetLastOperationRequest
	Res  func() *framework.GetLastOperationResponse
}

// GetLastOperation is the LastOperationGetter implementation. It just records the
// GetLastOperationRequest and returns a hardcoded GetLastOperationResponse.
func (l *LastOperationGetter) GetLastOperation(
	ctx context.Context,
	req *framework.GetLastOperationRequest,
) (*framework.GetLastOperationResponse, error) {

	l.Reqs = append(l.Reqs, req)
	return l.Res(), nil
}

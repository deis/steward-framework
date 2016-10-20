package fake

import (
	"context"

	"github.com/deis/steward-framework"
)

// LastOperationGetter is a fake implementation of a framework.LastOperationGetter. It's useful for
// use in unit tests
type LastOperationGetter struct {
	Reqs []*framework.GetLastOperationRequest
	Ret  func() *framework.GetLastOperationResponse
}

// GetLastOperation appends to l.Calls and returns l.Ret(), nil. Not concurrency safe
func (l *LastOperationGetter) GetLastOperation(
	ctx context.Context,
	req *framework.GetLastOperationRequest,
) (*framework.GetLastOperationResponse, error) {

	l.Reqs = append(l.Reqs, req)
	return l.Ret(), nil
}

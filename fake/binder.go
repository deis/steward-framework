package fake

import (
	"context"

	"github.com/deis/steward-framework"
)

// Binder is a fake (github.com/deis/steward-framework).Binder implementation, suitable for use in
// unit tests
type Binder struct {
	Reqs []*framework.BindRequest
	Res  *framework.BindResponse
}

// Bind is the Binder interface implementation. It constructs a new BindCall from the function
// params, then returns b.Res, nil. This function is not concurrency safe
func (b *Binder) Bind(
	ctx context.Context,
	req *framework.BindRequest,
) (*framework.BindResponse, error) {

	b.Reqs = append(b.Reqs, req)
	return b.Res, nil
}

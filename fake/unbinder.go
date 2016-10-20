package fake

import (
	"context"

	"github.com/deis/steward-framework"
)

// Unbinder is a fake (github.com/deis/steward-framework).Unbinder implementation, suitable for use
// in unit tests
type Unbinder struct {
	Reqs []*framework.UnbindRequest
}

// Unbind is the Unbinder interface implementaion. It returns nil
func (u *Unbinder) Unbind(ctx context.Context, req *framework.UnbindRequest) error {
	u.Reqs = append(u.Reqs, req)
	return nil
}

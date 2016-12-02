package fake

import (
	"context"

	"github.com/deis/steward-framework"
)

// Unbinder is a fake implementation of the (github.com/deis/steward-framework).Unbinder interface
// that is suitable for use in unit tests.
type Unbinder struct {
	Reqs []*framework.UnbindRequest
}

// Unbind is the Unbinder interface implementation. It just records the UnbindRequest and returns
// nil.
func (u *Unbinder) Unbind(
	ctx context.Context,
	serviceBrokerSpec framework.ServiceBrokerSpec,
	req *framework.UnbindRequest,
) error {
	u.Reqs = append(u.Reqs, req)
	return nil
}

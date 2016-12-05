package fake

import (
	"context"

	"github.com/deis/steward-framework"
)

// Binder is a fake implementation of the (github.com/deis/steward-framework).Binder interface
// that is suitable for use in unit tests.
type Binder struct {
	Reqs []*framework.BindRequest
	Res  *framework.BindResponse
}

// Bind is the Binder interface implementation. It just records the BindRequest and returns a
// hardcoded BindResponse.
func (b *Binder) Bind(
	ctx context.Context,
	serviceBrokerSpec framework.ServiceBrokerSpec,
	req *framework.BindRequest,
) (*framework.BindResponse, error) {

	b.Reqs = append(b.Reqs, req)
	return b.Res, nil
}

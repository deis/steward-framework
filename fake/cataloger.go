package fake

import (
	"context"

	"github.com/deis/steward-framework"
)

// Cataloger is a fake implementation of the (github.com/deis/steward-framework).Cataloger
// interface that is suitable for use in unit tests.
type Cataloger struct {
	Services []*framework.Service
}

// List is the Cataloger interface implementation. It returns a hardcoded, empty array of Service
// pointers.
func (f Cataloger) List(ctx context.Context, spec framework.ServiceBrokerSpec) ([]*framework.Service, error) {
	return f.Services, nil
}

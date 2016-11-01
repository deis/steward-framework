package fake

import (
	"github.com/deis/steward-framework"
)

// Cataloger is a fake implementation of the (github.com/deis/steward-framework).Cataloger
// interface that is suitable for use in unit tests.
type Cataloger struct {
	Services []*framework.Service
}

// List is the Cataloger interface implementation. It returns a hardcoded, empty array of Service
// pointers.
func (f Cataloger) List() ([]*framework.Service, error) {
	return f.Services, nil
}

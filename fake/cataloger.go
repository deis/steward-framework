package fake

import (
	"github.com/deis/steward-framework"
)

// Cataloger is a fake (github.com/deis/steward-framework).Cataloger implementation, suitable for use in unit tests
type Cataloger struct {
	Services []*framework.Service
}

// List is the Cataloger interface implementation. It returns f.Services
func (f Cataloger) List() ([]*framework.Service, error) {
	return f.Services, nil
}

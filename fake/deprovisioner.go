package fake

import (
	"context"

	"github.com/deis/steward-framework"
)

// Deprovisioner is a fake implementation of the (github.com/deis/steward-framework).Deprovisioner
// interface that is suitable for use in unit tests.
type Deprovisioner struct {
	Reqs []*framework.DeprovisionRequest
	Res  *framework.DeprovisionResponse
	Err  error
}

// Deprovision is the Deprovisioner interface implementation. It just records the
// DeprovisionRequest and returns a hardcoded DeprovisionResponse.
func (d *Deprovisioner) Deprovision(
	ctx context.Context,
	serviceBrokerSpec framework.ServiceBrokerSpec,
	req *framework.DeprovisionRequest,
) (*framework.DeprovisionResponse, error) {

	d.Reqs = append(d.Reqs, req)
	return d.Res, d.Err
}

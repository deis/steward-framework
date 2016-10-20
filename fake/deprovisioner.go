package fake

import (
	"context"

	"github.com/deis/steward-framework"
)

// Deprovisioner is a fake implementation of (github.com/deis/steward-framework).Deprovisioner,
// suitable for usage in unit tests
type Deprovisioner struct {
	Reqs []*framework.DeprovisionRequest
	Resp *framework.DeprovisionResponse
	Err  error
}

// Deprovision is the Deprovisioner interface implementation. It packages the function parameters
// into a DeprovisionCall, appends it to d.Deprovisons, and returns d.Resp, d.Err. This function is not concurrency safe
func (d *Deprovisioner) Deprovision(
	ctx context.Context,
	req *framework.DeprovisionRequest,
) (*framework.DeprovisionResponse, error) {

	d.Reqs = append(d.Reqs, req)
	return d.Resp, d.Err
}

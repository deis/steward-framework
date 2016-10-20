package fake

import (
	"context"

	"github.com/deis/steward-framework"
)

// Provisioner is a fake implementation of (github.com/deis/steward-framework).Provisioner, suitable for usage in unit tests
type Provisioner struct {
	Reqs []*framework.ProvisionRequest
	Resp *framework.ProvisionResponse
	Err  error
}

// Provision is the Provisioner interface implementation. It packages the function parameters into a ProvisionCall, stores them in p.Provisioned, and returns p.Resp, p.Err. This function is not concurrency safe
func (p *Provisioner) Provision(
	ctx context.Context,
	req *framework.ProvisionRequest,
) (*framework.ProvisionResponse, error) {

	p.Reqs = append(p.Reqs, req)
	return p.Resp, p.Err
}

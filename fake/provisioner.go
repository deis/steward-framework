package fake

import (
	"context"

	"github.com/deis/steward-framework"
)

// Provisioner is a fake implementation of the (github.com/deis/steward-framework).Provisioner
// interface that is suitable for use in unit tests.
type Provisioner struct {
	Reqs []*framework.ProvisionRequest
	Res  *framework.ProvisionResponse
	Err  error
}

// Provision is the Provisioner interface implementation. It just records the ProvisionRequest and
// returns a hardcoded ProvisionResponse.
func (p *Provisioner) Provision(
	ctx context.Context,
	serviceBrokerSpec framework.ServiceBrokerSpec,
	req *framework.ProvisionRequest,
) (*framework.ProvisionResponse, error) {

	p.Reqs = append(p.Reqs, req)
	return p.Res, p.Err
}

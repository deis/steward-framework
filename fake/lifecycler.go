package fake

import (
	"github.com/deis/steward-framework"
)

// Lifecycler is a composition of the provisioner, deprovisioner, binder and unbinder. It's intended for use in passing to functions that require all functionality
type Lifecycler struct {
	framework.Provisioner
	framework.Deprovisioner
	framework.Binder
	framework.Unbinder
	framework.LastOperationGetter
}

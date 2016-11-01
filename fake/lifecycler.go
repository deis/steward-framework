package fake

import (
	"github.com/deis/steward-framework"
)

// Lifecycler is a composition of this package's fake implementations of the Provisioner,
// Deprovisioner, Binder, Unbinder, and LastOperationGetter interfaces. It is suitable for use
// in unit tests.
type Lifecycler struct {
	framework.Provisioner
	framework.Deprovisioner
	framework.Binder
	framework.Unbinder
	framework.LastOperationGetter
}

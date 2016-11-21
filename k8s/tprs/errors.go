package tprs

import (
	"fmt"
)

type errCreatingThirdPartyResource struct {
	Original error
}

func (e errCreatingThirdPartyResource) Error() string {
	return fmt.Sprintf("error creating third party resource: %s", e.Original)
}

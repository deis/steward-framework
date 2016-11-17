// +build integration

package clients

import (
	"testing"

	"github.com/arschles/assert"
	testk8s "github.com/deis/steward-framework/testing/k8s"
)

func TestNewDynamicClient(t *testing.T) {
	restCfg, err := testk8s.GetRESTConfig()
	assert.NoErr(t, err)
	dynCl, err := NewDynamic(*restCfg)
	assert.NotNil(t, dynCl, "dynamic client")
	assert.NoErr(t, err)
}

// +build integration

package binding

import (
	"testing"

	"github.com/arschles/assert"
	"github.com/deis/steward-framework/k8s"
	testk8s "github.com/deis/steward-framework/testing/k8s"
)

func TestNewK8sWatchBindingFunc(t *testing.T) {
	restCfg, err := testk8s.NewRESTConfig()
	assert.NoErr(t, err)
	dynCl, err := k8s.NewDynamicClient(*restCfg)
	assert.NoErr(t, err)
	fn := NewK8sWatchBindingFunc(dynCl)
	_, err := fn("default")
	assert.NoErr(t, err)
	// TODO: do the watch
}

package restutil

import (
	"strings"
	"testing"

	"github.com/arschles/assert"
)

func TestAbsPath(t *testing.T) {
	const (
		apiVersionBase = "testbase"
		apiVersion     = "testversion"
		namespace      = "testNS"
		pluralName     = "testplural"
	)
	elts := AbsPath(apiVersionBase, apiVersion, true, namespace, pluralName)
	assert.Equal(t, len(elts), 7, "number of path elts")
	assert.Equal(t, elts[0], "apis", "first path elt")
	assert.Equal(t, elts[1], apiVersionBase, "base API version")
	assert.Equal(t, elts[2], apiVersion, "api version")
	assert.Equal(t, elts[3], "watch", "watch path element")
	assert.Equal(t, elts[4], "namespaces", "'namespaces' path elt")
	assert.Equal(t, elts[5], namespace, "namespace value")
	assert.Equal(t, elts[6], strings.ToLower(pluralName), "plural name")

	elts = AbsPath(apiVersionBase, apiVersion, false, namespace, pluralName)
	assert.Equal(t, len(elts), 6, "number of path elts")
}

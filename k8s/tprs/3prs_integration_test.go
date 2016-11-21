// +build integration

package tprs

import (
	"testing"

	"github.com/arschles/assert"
)

func TestEnsure3PRs(t *testing.T) {
	err := Ensure3PRs(clientset)
	assert.NoErr(t, err)
}

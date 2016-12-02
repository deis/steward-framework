// +build integration

package binding

import (
	"testing"

	"github.com/deis/steward-framework/k8s/tprs"
	"github.com/deis/steward-framework/testing/k8s"
	testsetup "github.com/deis/steward-framework/testing/setup"
	"github.com/technosophos/moniker"
)

var (
	testNamespace string
)

func TestMain(m *testing.M) {
	testsetup.SetupAndTearDown(m, setup, teardown)
}

func setup() error {
	k8sClient, err := k8s.GetClientset()
	if err != nil {
		return err
	}
	if err := tprs.Ensure3PRs(k8sClient); err != nil {
		return err
	}
	testNamespace = moniker.New().NameSep("-")
	if err := k8s.EnsureNamespace(testNamespace); err != nil {
		return err
	}
	return nil
}

func teardown() error {
	// This will also delete the serviceBroker
	if err := k8s.DeleteNamespace(testNamespace); err != nil {
		return err
	}
	return nil
}

// +build integration

package k8s

import (
	"testing"

	"github.com/deis/steward-framework/testing/k8s"
	testsetup "github.com/deis/steward-framework/testing/setup"
	"k8s.io/client-go/kubernetes"
)

var (
	clientset *kubernetes.Clientset
)

func TestMain(m *testing.M) {
	testsetup.SetupAndTearDown(m, setup, teardown)
}

func setup() error {
	var err error
	if clientset, err = k8s.GetClientset(); err != nil {
		return err
	}
	return nil
}

func teardown() error {
	return nil
}

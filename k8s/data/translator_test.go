package data

import (
	"testing"

	"github.com/arschles/assert"
	"k8s.io/client-go/pkg/api"
)

func TestTranslateToTPRMismatchedTypes(t *testing.T) {
	raw := new(api.Pod)
	raw.Kind = "Pod"
	serviceBinding := new(ServiceBinding)
	err := TranslateToTPR(raw, serviceBinding, ServiceBindingKind)
	assert.NotNil(t, err, "returned error")
	mismatchErr, ok := err.(ErrMismatchedKinds)
	assert.True(t, ok, "returned error was not a ErrMismatchedKinds")
	assert.Equal(t, mismatchErr.RawKind, raw.Kind, "raw kind")
	assert.Equal(t, mismatchErr.Expected, ServiceBindingKind, "expected kind")
}

func TestTranslateToTPRMatchingTypes(t *testing.T) {
	raw := new(ServiceBinding)
	raw.Kind = "ServiceBinding"
	raw.Name = "myservicebinding"
	raw.Namespace = "myns"
	raw.Spec.ID = "myid"
	serviceBinding := new(ServiceBinding)
	assert.NoErr(t, TranslateToTPR(raw, serviceBinding, ServiceBindingKind))
	assert.NotNil(t, serviceBinding, "service binding")
	assert.Equal(t, raw.Kind, serviceBinding.Kind, "kind")
	assert.Equal(t, raw.Name, serviceBinding.Name, "name")
	assert.Equal(t, raw.Namespace, serviceBinding.Namespace, "namespace")
	assert.Equal(t, raw.Spec.ID, serviceBinding.Spec.ID, "ID")
}

func TestTranslateToUnstructued(t *testing.T) {
	serviceBinding := new(ServiceBinding)
	serviceBinding.Kind = ServiceBindingKind
	serviceBinding.Namespace = "myns"
	serviceBinding.Name = "myname"
	unstruc, err := TranslateToUnstructured(serviceBinding)
	assert.NoErr(t, err)
	assert.Equal(t, unstruc.GetKind(), serviceBinding.Kind, "kind")
	assert.Equal(t, unstruc.GetNamespace(), serviceBinding.Namespace, "namespace")
	assert.Equal(t, unstruc.GetName(), serviceBinding.Name, "name")
}

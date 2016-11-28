package data

import (
	"testing"

	"github.com/arschles/assert"
	"k8s.io/client-go/pkg/api"
)

func TestTranslateToTPRMismatchedTypes(t *testing.T) {
	raw := new(api.Pod)
	raw.Kind = "Pod"
	binding := new(Binding)
	err := TranslateToTPR(raw, binding, BindingKind)
	assert.NotNil(t, err, "returned error")
	mismatchErr, ok := err.(ErrMismatchedKinds)
	assert.True(t, ok, "returned error was not a ErrMismatchedKinds")
	assert.Equal(t, mismatchErr.RawKind, raw.Kind, "raw kind")
	assert.Equal(t, mismatchErr.Expected, BindingKind, "expected kind")
}

func TestTranslateToTPRMatchingTypes(t *testing.T) {
	raw := new(Binding)
	raw.Kind = "Binding"
	raw.Name = "mybinding"
	raw.Namespace = "myns"
	raw.Spec.ID = "myid"
	binding := new(Binding)
	assert.NoErr(t, TranslateToTPR(raw, binding, BindingKind))
	assert.NotNil(t, binding, "binding")
	assert.Equal(t, raw.Kind, binding.Kind, "kind")
	assert.Equal(t, raw.Name, binding.Name, "name")
	assert.Equal(t, raw.Namespace, binding.Namespace, "namespace")
	assert.Equal(t, raw.Spec.ID, binding.Spec.ID, "ID")
}

func TestTranslateToUnstructued(t *testing.T) {
	binding := new(Binding)
	binding.Kind = BindingKind
	binding.Namespace = "myns"
	binding.Name = "myname"
	unstruc, err := TranslateToUnstructured(binding)
	assert.NoErr(t, err)
	assert.Equal(t, unstruc.GetKind(), binding.Kind, "kind")
	assert.Equal(t, unstruc.GetNamespace(), binding.Namespace, "namespace")
	assert.Equal(t, unstruc.GetName(), binding.Name, "name")
}

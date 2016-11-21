// +build integration

package binding

import (
	"testing"
	"time"

	"github.com/arschles/assert"
	"github.com/deis/steward-framework/k8s/clients"
	"github.com/deis/steward-framework/k8s/data"
	testk8s "github.com/deis/steward-framework/testing/k8s"
	"k8s.io/client-go/pkg/api"
	"k8s.io/client-go/pkg/api/unversioned"
	"k8s.io/client-go/pkg/watch"
)

func TestNewK8sWatchBindingFunc(t *testing.T) {
	const (
		ns          = "test"
		bindingName = "testBinding"
		timeout     = 500 * time.Millisecond
	)
	assert.NoErr(t, testk8s.EnsureNamespace(ns))
	restCfg, err := testk8s.GetRESTConfig()
	assert.NoErr(t, err)
	dynCl, err := clients.NewDynamic(*restCfg)
	assert.NoErr(t, err)
	fn := NewK8sWatchBindingFunc(dynCl)
	watcher, err := fn(ns)
	assert.NoErr(t, err)
	defer watcher.Stop()
	ch := watcher.ResultChan()
	binding := data.Binding{
		TypeMeta: unversioned.TypeMeta{
			Kind:       data.BindingKind,
			APIVersion: data.APIVersion,
		},
		ObjectMeta: api.ObjectMeta{
			Name:      bindingName,
			Namespace: ns,
		},
		Spec:   data.BindingSpec{},
		Status: data.BindingStatus{},
	}
	unstructuredBinding, err := data.TranslateToUnstructured(&binding)
	assert.NoErr(t, err)
	resourceCl := dynCl.Resource(data.BindingAPIResource(), ns)
	_, createErr := resourceCl.Create(unstructuredBinding)
	assert.NoErr(t, createErr)
	select {
	case evt := <-ch:
		assert.Equal(t, evt.Type, watch.Added, "type")
		binding, ok := evt.Object.(*data.Binding)
		assert.True(t, ok, "returned event object was not a binding")
		assert.Equal(t, binding.Kind, data.BindingKind, "kind")
		assert.Equal(t, binding.APIVersion, data.APIVersion, "api version")
		assert.Equal(t, binding.Name, bindingName, "name")
		assert.Equal(t, binding.Namespace, ns, "namespace")
	case <-time.After(timeout):
		t.Fatalf("didn't receive an event within %s", timeout)
	}

}

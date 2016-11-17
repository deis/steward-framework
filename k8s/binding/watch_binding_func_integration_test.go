// +build integration

package binding

import (
	"testing"
	"time"

	"github.com/arschles/assert"
	"github.com/deis/steward-framework/k8s"
	"github.com/deis/steward-framework/k8s/data"
	testk8s "github.com/deis/steward-framework/testing/k8s"
	"k8s.io/client-go/watch"
)

func TestNewK8sWatchBindingFunc(t *testing.T) {
	const (
		ns          = "test"
		bindingName = "testBinding"
		timeout     = 500 * time.Millisecond
	)
	assert.NoErr(t, testk8s.EnsureNamespace(ns))
	restCfg, err := testk8s.NewRESTConfig()
	assert.NoErr(t, err)
	dynCl, err := k8s.NewDynamicClient(*restCfg)
	assert.NoErr(t, err)
	fn := NewK8sWatchBindingFunc(dynCl)
	watcher, err := fn(ns)
	assert.NoErr(t, err)
	defer watcher.Stop()
	ch := watcher.ResultCh()
	unstructuredBinding := runtime.Unstructured{
		Object: map[string]interface{}{
			"Kind":       data.BindingKind,
			"APIVersion": "steward.deis.io/v1",
			"Metadata": map[string]string{
				"Name":      bindingName,
				"Namespace": ns,
			},
			"Spec":   map[string]string{},
			"Status": map[string]string{},
		},
	}
	resourceCl := dynCl.Resource(&bindingAPIResource, ns)
	_, createErr := resourceCl.Create(&unstructuredBinding)
	assert.NoErr(t, createErr)
	select {
	case evt := <-ch:
		assert.Equal(t, evt.Type, watch.Added)
		binding, ok := evt.Object.(*data.Binding)
		assert.True(t, ok, "returned event object was not a binding")
		asssert.Equal(t, binding.Kind, data.BindingKind, "kind")
		assert.Equal(t, binding.APIVersion, "steward.deis.io/v1", "api version")
		assert.Equal(t, binding.Name, bindingName, "name")
		assert.Equal(t, binding.Namespace, ns, "namespace")
	case <-time.After(timeout):
		t.Fatalf("didn't receive an event within %s", timeout)
	}

}

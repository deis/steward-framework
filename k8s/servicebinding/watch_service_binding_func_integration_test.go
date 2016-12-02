// +build integration

package servicebinding

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

func TestNewK8sWatchServiceBindingFunc(t *testing.T) {
	const (
		serviceBindingName = "test-service-binding"
		timeout            = 500 * time.Millisecond
	)
	restCfg, err := testk8s.GetRESTConfig()
	assert.NoErr(t, err)
	dynCl, err := clients.NewDynamic(*restCfg)
	assert.NoErr(t, err)
	fn := NewK8sWatchServiceBindingFunc(dynCl)
	watcher, err := fn(testNamespace)
	assert.NoErr(t, err)
	defer watcher.Stop()
	ch := watcher.ResultChan()
	origServiceBinding := data.ServiceBinding{
		TypeMeta: unversioned.TypeMeta{
			Kind:       data.ServiceBindingKind,
			APIVersion: data.APIVersion,
		},
		ObjectMeta: api.ObjectMeta{
			Name:      serviceBindingName,
			Namespace: testNamespace,
		},
		Spec:   data.ServiceBindingSpec{},
		Status: data.ServiceBindingStatus{},
	}
	unstructuredServiceBinding, err := data.TranslateToUnstructured(&origServiceBinding)
	assert.NoErr(t, err)
	resourceCl := dynCl.Resource(data.ServiceBindingAPIResource(), testNamespace)
	_, createErr := resourceCl.Create(unstructuredServiceBinding)
	assert.NoErr(t, createErr)
	select {
	case evt := <-ch:
		assert.Equal(t, evt.Type, watch.Added, "type")
		serviceBinding := new(data.ServiceBinding)
		assert.NoErr(t, data.TranslateToTPR(evt.Object, serviceBinding, data.ServiceBindingKind))
		assert.Equal(t, serviceBinding.Kind, data.ServiceBindingKind, "kind")
		assert.Equal(t, serviceBinding.APIVersion, data.APIVersion, "api version")
		assert.Equal(t, serviceBinding.Name, serviceBindingName, "name")
		assert.Equal(t, serviceBinding.Namespace, testNamespace, "namespace")
	case <-time.After(timeout):
		t.Fatalf("didn't receive an event within %s", timeout)
	}

}

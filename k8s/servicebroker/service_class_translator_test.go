package servicebroker

import (
	"fmt"
	"testing"

	"github.com/arschles/assert"
	"github.com/deis/steward-framework"
	"github.com/deis/steward-framework/k8s/data"
	"k8s.io/client-go/pkg/api/unversioned"
	"k8s.io/client-go/pkg/api/v1"
)

func TestTranslateServiceClass(t *testing.T) {
	parentServiceBroker := &data.ServiceBroker{
		TypeMeta: unversioned.TypeMeta{},
		ObjectMeta: v1.ObjectMeta{
			Name:      "testName",
			Namespace: "testNS",
		},
	}
	svc := &framework.Service{
		ServiceInfo: framework.ServiceInfo{ID: "testid"},
	}
	sClass := translateServiceClass(parentServiceBroker, svc)
	assert.Equal(t, sClass.Kind, data.ServiceClassKind, "kind")
	assert.Equal(t, sClass.Name, serviceClassName(parentServiceBroker, svc), "name")
	assert.Equal(t, sClass.ID, svc.ID, "ID")
	assert.Equal(t, sClass.Namespace, parentServiceBroker.Namespace, "namespace")
	assert.Equal(t, len(sClass.Plans), len(svc.Plans), "number of plans")
}

func TestServiceClassName(t *testing.T) {
	serviceBroker := &data.ServiceBroker{
		ObjectMeta: v1.ObjectMeta{Name: "testServiceBroker"},
	}
	svc := &framework.Service{
		ServiceInfo: framework.ServiceInfo{Name: "testSvc"},
	}
	name := serviceClassName(serviceBroker, svc)
	assert.Equal(t, name, fmt.Sprintf("%s-%s", serviceBroker.Name, svc.Name), "service class name")
}

func TestTranslatePlans(t *testing.T) {
	plans := []framework.ServicePlan{
		{ID: "ID1", Name: "name1", Description: "descr1", Free: true},
		{ID: "ID2", Name: "name2", Description: "descr2", Free: false},
		{ID: "ID3", Name: "name3", Description: "descr3", Free: true},
	}
	translated := translatePlans(plans)
	assert.Equal(t, len(translated), len(plans), "number of translated plans")
	for i, plan := range plans {
		tr := translated[i]
		assert.Equal(t, tr.ID, plan.ID, "ID")
		assert.Equal(t, tr.Name, plan.Name, "name")
		assert.Equal(t, tr.Description, plan.Description, "description")
	}
}

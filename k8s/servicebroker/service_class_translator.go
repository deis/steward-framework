package servicebroker

import (
	"fmt"

	"github.com/deis/steward-framework"
	"github.com/deis/steward-framework/k8s/data"
	"k8s.io/client-go/pkg/api/unversioned"
	"k8s.io/client-go/pkg/api/v1"
)

func translateServiceClass(
	parentServiceBroker *data.ServiceBroker,
	svc *framework.Service) *data.ServiceClass {

	serviceBrokerRef := getObjectReference(parentServiceBroker)
	return &data.ServiceClass{
		TypeMeta: unversioned.TypeMeta{
			APIVersion: data.APIVersion,
			Kind:       data.ServiceClassKind,
		},
		ObjectMeta: v1.ObjectMeta{
			Name:      serviceClassName(parentServiceBroker, svc),
			Namespace: parentServiceBroker.Namespace,
		},
		ServiceBrokerRef:  *serviceBrokerRef,
		ID:                svc.ID,
		ServiceBrokerName: parentServiceBroker.Name,
		Bindable:          true,
		Plans:             translatePlans(svc.Plans),
		PlanUpdatable:     svc.PlanUpdatable,
		Description:       svc.Description,
	}
}

func serviceClassName(parentServiceBroker *data.ServiceBroker, svc *framework.Service) string {
	return fmt.Sprintf("%s-%s", parentServiceBroker.Name, svc.Name)
}

func translatePlans(plans []framework.ServicePlan) []data.ServicePlan {
	ret := make([]data.ServicePlan, len(plans))
	for i, plan := range plans {
		ret[i] = data.ServicePlan{ID: plan.ID, Name: plan.Name, Description: plan.Description}
	}
	return ret
}

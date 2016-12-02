package data

import (
	"strings"

	"k8s.io/client-go/pkg/api"
	"k8s.io/client-go/pkg/api/unversioned"
	"k8s.io/client-go/pkg/api/v1"
)

const (
	ServiceClassKind       = "ServiceClass"
	ServiceClassKindPlural = "ServiceClasses"
)

func ServiceClassAPIResource() *unversioned.APIResource {
	return &unversioned.APIResource{
		Name:       strings.ToLower(ServiceClassKindPlural),
		Namespaced: true,
		Kind:       ServiceClassKind,
	}
}

type ServiceClass struct {
	unversioned.TypeMeta `json:",inline"`
	v1.ObjectMeta        `json:"metadata,omitempty"`
	ServiceBrokerRef     api.ObjectReference `json:"service_broker_ref"`
	ID                   string              `json:"id"`
	ServiceBrokerName    string              `json:"service_broker_name"`
	Bindable             bool                `json:"bindable"`
	Plans                []ServicePlan       `json:"plans"`
	PlanUpdatable        bool                `json:"updatable"`
	Description          string              `json:"description"`
}

type ServicePlan struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

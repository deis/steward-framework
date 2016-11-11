package broker

import (
	"github.com/deis/steward-framework/k8s/data"
	"k8s.io/client-go/pkg/api"
)

func getObjectReference(br *data.Broker) *api.ObjectReference {
	return &api.ObjectReference{
		Kind:            data.BrokerKind,
		Namespace:       br.Namespace,
		Name:            br.Name,
		UID:             br.UID,
		APIVersion:      br.APIVersion,
		ResourceVersion: br.ResourceVersion,
	}
}

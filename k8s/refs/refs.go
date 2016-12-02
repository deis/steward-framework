package refs

import (
	"github.com/deis/steward-framework/k8s/data"
)

// GetDependenciesForInstance fetches the entire reference tree for inst
func GetDependenciesForInstance(
	inst *data.Instance,
	getServiceBroker ServiceBrokerGetterFunc,
	getServiceClass ServiceClassGetterFunc,
) (*data.ServiceBroker, *data.ServiceClass, error) {
	sClass, err := getServiceClass(inst.Spec.ServiceClassRef)
	if err != nil {
		return nil, nil, err
	}
	serviceBroker, err := getServiceBroker(sClass.ServiceBrokerRef)
	if err != nil {
		return nil, nil, err
	}
	return serviceBroker, sClass, nil
}

// GetDependenciesForBinding fetches the entire reference tree for binding
func GetDependenciesForBinding(
	binding *data.Binding,
	getServiceBroker ServiceBrokerGetterFunc,
	getServiceClass ServiceClassGetterFunc,
	getInstance InstanceGetterFunc,
) (*data.ServiceBroker, *data.ServiceClass, *data.Instance, error) {
	inst, err := getInstance(binding.Spec.InstanceRef)
	if err != nil {
		return nil, nil, nil, err
	}
	serviceBroker, sClass, err := GetDependenciesForInstance(inst, getServiceBroker, getServiceClass)
	if err != nil {
		return nil, nil, nil, err
	}
	return serviceBroker, sClass, inst, nil
}

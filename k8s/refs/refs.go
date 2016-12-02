package refs

import (
	"github.com/deis/steward-framework/k8s/data"
)

// GetDependenciesForServiceInstance fetches the entire reference tree for inst
func GetDependenciesForServiceInstance(
	inst *data.ServiceInstance,
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
	getServiceInstance ServiceInstanceGetterFunc,
) (*data.ServiceBroker, *data.ServiceClass, *data.ServiceInstance, error) {
	inst, err := getServiceInstance(binding.Spec.ServiceInstanceRef)
	if err != nil {
		return nil, nil, nil, err
	}
	serviceBroker, sClass, err := GetDependenciesForServiceInstance(inst, getServiceBroker, getServiceClass)
	if err != nil {
		return nil, nil, nil, err
	}
	return serviceBroker, sClass, inst, nil
}

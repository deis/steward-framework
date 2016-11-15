package refs

import (
	"github.com/deis/steward-framework/k8s/data"
)

// GetDependenciesForInstance fetches the entire reference tree for inst
func GetDependenciesForInstance(
	inst *data.Instance,
	getBroker BrokerGetterFunc,
	getServiceClass ServiceClassGetterFunc,
) (*data.Broker, *data.ServiceClass, error) {
	sClass, err := getServiceClass(inst.Spec.ServiceClassRef)
	if err != nil {
		return nil, nil, err
	}
	broker, err := getBroker(sClass.BrokerRef)
	if err != nil {
		return nil, nil, err
	}
	return broker, sClass, nil
}

// GetDependenciesForBinding fetches the entire reference tree for binding
func GetDependenciesForBinding(
	binding *data.Binding,
	getBroker BrokerGetterFunc,
	getServiceClass ServiceClassGetterFunc,
	getInstance InstanceGetterFunc,
) (*data.Broker, *data.ServiceClass, *data.Instance, error) {
	inst, err := getInstance(binding.Spec.InstanceRef)
	if err != nil {
		return nil, nil, nil, err
	}
	broker, sClass, err := GetDependenciesForInstance(inst, getBroker, getServiceClass)
	if err != nil {
		return nil, nil, nil, err
	}
	return broker, sClass, inst, nil
}

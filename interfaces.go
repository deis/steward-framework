package framework

import (
	"context"
)

// Cataloger lists all the available services.
type Cataloger interface {
	// List lists all the available services.
	List(
		ctx context.Context,
		serviceBrokerSpec ServiceBrokerSpec,
	) ([]*Service, error)
}

// Provisioner provisions services instances.
type Provisioner interface {
	// Provision provisions a service instance.
	Provision(
		ctx context.Context,
		serviceBrokerSpec ServiceBrokerSpec,
		req *ProvisionRequest,
	) (*ProvisionResponse, error)
}

// Binder obtains valid credentials (and other connection details) for service instances.
type Binder interface {
	// Bind obtains valid credentials (and other connection details) for a service instance.
	Bind(
		ctx context.Context,
		serviceBrokerSpec ServiceBrokerSpec,
		req *BindRequest,
	) (*BindResponse, error)
}

// Unbinder invalidates credentials for service instances.
type Unbinder interface {
	// Unbind invalidates the specified credentials.
	Unbind(
		ctx context.Context,
		serviceBrokerSpec ServiceBrokerSpec,
		req *UnbindRequest,
	) error
}

// Deprovisioner deprovisions service instances.
type Deprovisioner interface {
	// Deprovision deprovisions a service instance.
	Deprovision(
		ctx context.Context,
		serviceBrokerSpec ServiceBrokerSpec,
		req *DeprovisionRequest,
	) (*DeprovisionResponse, error)
}

// OperationStatusRetriever fetches the status of an asynchronous operation that is pending
// completion.
type OperationStatusRetriever interface {
	// GetOperationStatus fetches the status of an asynchronous operation that is pending completion.
	GetOperationStatus(
		ctx context.Context,
		serviceBrokerSpec ServiceBrokerSpec,
		req *OperationStatusRequest,
	) (*OperationStatusResponse, error)
}

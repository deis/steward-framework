package framework

import (
	"context"
)

// Cataloger lists all the available services
type Cataloger interface {
	List(ctx context.Context) ([]*Service, error)
}

// Provisioner provisions services
type Provisioner interface {
	Provision(ctx context.Context, req *ProvisionRequest) (*ProvisionResponse, error)
}

// Binder binds services to apps
type Binder interface {
	Bind(ctx context.Context, req *BindRequest) (*BindResponse, error)
}

// Unbinder unbinds services from apps
type Unbinder interface {
	Unbind(ctx context.Context, req *UnbindRequest) error
}

// Deprovisioner deprovisions services
type Deprovisioner interface {
	Deprovision(ctx context.Context, req *DeprovisionRequest) (*DeprovisionResponse, error)
}

// LastOperationGetter fetches the last operation performed after an async provision or deprovision response
type LastOperationGetter interface {
	GetLastOperation(ctx context.Context, req *GetLastOperationRequest) (*GetLastOperationResponse, error)
}

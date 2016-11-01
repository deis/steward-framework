package framework

// BindRequest represents a request to bind to a service.
type BindRequest struct {
	InstanceID string
	ServiceID  string
	PlanID     string
	BindingID  string
	Parameters map[string]interface{}
}

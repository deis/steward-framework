package framework

// UnbindRequest represents a request to unbind from a service.
type UnbindRequest struct {
	InstanceID string
	ServiceID  string
	PlanID     string
	BindingID  string
}

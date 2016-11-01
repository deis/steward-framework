package framework

// UnbindRequest represents a request invalidate a consumer's credentials for a service instance.
type UnbindRequest struct {
	InstanceID string
	ServiceID  string
	PlanID     string
	BindingID  string
}

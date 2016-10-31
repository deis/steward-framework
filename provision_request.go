package framework

// ProvisionRequest represents a request to do a service provision operation.
type ProvisionRequest struct {
	InstanceID        string
	PlanID            string
	ServiceID         string
	AcceptsIncomplete bool
	Parameters        map[string]interface{}
}

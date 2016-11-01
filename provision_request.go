package framework

// ProvisionRequest represents a request to provision a service instance.
type ProvisionRequest struct {
	InstanceID        string
	PlanID            string
	ServiceID         string
	AcceptsIncomplete bool
	Parameters        map[string]interface{}
}

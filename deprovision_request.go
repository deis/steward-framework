package framework

// DeprovisionRequest represents a request to deprovision a service instance.
type DeprovisionRequest struct {
	InstanceID        string
	ServiceID         string
	PlanID            string
	AcceptsIncomplete bool
	Parameters        map[string]interface{}
}

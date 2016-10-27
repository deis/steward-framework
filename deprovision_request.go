package framework

// DeprovisionRequest represents a request to do a service deprovision operation.
type DeprovisionRequest struct {
	InstanceID        string
	ServiceID         string
	PlanID            string
	AcceptsIncomplete bool
	Parameters        map[string]interface{}
}

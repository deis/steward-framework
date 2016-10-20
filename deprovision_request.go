package framework

// DeprovisionRequest represents a request to do a service deprovision operation. This struct is JSON-compatible with the request body detailed at https://docs.cloudfoundry.org/services/api.html#deprovisioning
type DeprovisionRequest struct {
	InstanceID        string     `json:"instance_id"`
	ServiceID         string     `json:"service_id"`
	PlanID            string     `json:"plan_id"`
	AcceptsIncomplete bool       `json:"accepts_incomplete"`
	Parameters        JSONObject `json:parameters,omitempty`
}

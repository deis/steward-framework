package framework

// GetLastOperationRequest represents a request to query the status of an asynchronous operation
// that is pending completion by the backing broker.
type GetLastOperationRequest struct {
	InstanceID string
	ServiceID  string
	PlanID     string
	Operation  string
}

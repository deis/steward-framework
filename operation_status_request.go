package framework

// OperationStatusRequest represents a request to query the status of an asynchronous operation
// that is pending completion by the backing service broker.
type OperationStatusRequest struct {
	InstanceID string
	ServiceID  string
	PlanID     string
	Operation  string
}

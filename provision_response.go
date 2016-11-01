package framework

// ProvisionResponse represents a response to ProvisionRequest.
type ProvisionResponse struct {
	Operation string
	IsAsync   bool
	Extra     map[string]interface{}
}

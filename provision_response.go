package framework

// ProvisionResponse represents a response to a provisioning request.
type ProvisionResponse struct {
	Operation string
	IsAsync   bool
	Extra     map[string]interface{}
}

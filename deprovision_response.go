package framework

// DeprovisionResponse represents a response to a DeprovisionRequest.
type DeprovisionResponse struct {
	Operation string
	IsAsync   bool
}

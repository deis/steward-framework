package framework

// BindResponse represents a response to a BindRequest. It contains a map of credentials and other
// connection details for a service instance.
type BindResponse struct {
	Creds map[string]interface{}
}

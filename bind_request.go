package framework

const (
	// TargetNamespaceKey is the required key for the target namespace value
	TargetNamespaceKey = "target_namespace"
	// TargetNameKey is the required key for the target name value
	TargetNameKey = "target_name"
)

// BindRequest represents a request to bind to a service.
type BindRequest struct {
	InstanceID string
	ServiceID  string
	PlanID     string
	BindingID  string
	Parameters map[string]interface{}
}

// TargetNamespace returns the target namespace in b.Parameters, or an error if it's missing
func (b BindRequest) TargetNamespace() (string, error) {
	return b.paramString(TargetNamespaceKey)
}

// TargetName returns the target name in b.Parameters, or an error if it's missing
func (b BindRequest) TargetName() (string, error) {
	return b.paramString(TargetNameKey)
}

func (b BindRequest) paramString(key string) (string, error) {
	i, ok := b.Parameters[key]
	if !ok {
		return "", errMissing
	}
	s, ok := i.(string)
	if !ok {
		return "", errNotAString{value: i}
	}
	return s, nil
}

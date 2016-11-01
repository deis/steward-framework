package framework

// ServicePlan represents a specific configuration or variant of its parent Service.
type ServicePlan struct {
	ID          string
	Name        string
	Description string
	Free        bool
}

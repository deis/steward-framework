package framework

// ServiceInfo represents all of the information about a service except for its plans
type ServiceInfo struct {
	Name          string
	ID            string
	Description   string
	PlanUpdatable bool
}

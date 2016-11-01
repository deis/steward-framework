package framework

// ServiceInfo represents all ascpects of a requestable service, except for its ServicePlans.
type ServiceInfo struct {
	Name          string
	ID            string
	Description   string
	PlanUpdatable bool
}

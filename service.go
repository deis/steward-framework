package framework

// Service represents a requestable service.
type Service struct {
	ServiceInfo
	Plans []ServicePlan
}

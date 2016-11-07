package framework

// Lifecycler is a composition of several more narrowly scoped interfaces. It's intended for use in
// passing to functions that require all functionality.
type Lifecycler interface {
	Provisioner
	Deprovisioner
	Binder
	Unbinder
	OperationStatusRetriever
}

package component

// Context holds dependencies for components
type Context map[string]Lifecycle

// Lifecycle interface for component lifecycle management
type Lifecycle interface {
	// Start initializes the component
	Start(ctx Context) (Lifecycle, error)

	// Stop shuts down the component
	Stop(ctx Context) error
}

package component

import (
	"fmt"
	"sync"
	"time"
)

// Component defines a component with its dependencies
type Component struct {
	key          string
	instance     Lifecycle
	dependencies []string
	result       interface{}
	started      bool
	mu           sync.Mutex
}

// Define creates a new component
func Define(key string, instance Lifecycle, dependencies ...string) *Component {
	return &Component{
		key:          key,
		instance:     instance,
		dependencies: dependencies,
		started:      false,
	}
}

func (c *Component) Key() string {
	return c.key
}

// Start initializes the component
func (c *Component) Start(ctx Context) (Lifecycle, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.started {
		return c.instance, nil
	}

	startTime := time.Now()
	result, err := c.instance.Start(ctx)
	elapsedTime := time.Since(startTime)
	
	fmt.Printf("Component %s started successfully in %v\n", c.key, elapsedTime)
	if err != nil {
		return nil, fmt.Errorf("failed to start component: %w", err)
	}

	c.result = result
	c.started = true
	return result, nil
}

// Stop shuts down the component
func (c *Component) Stop(ctx Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.started {
		return nil
	}

	err := c.instance.Stop(ctx)
	fmt.Printf("Component %s stopped successfully\n", c.key)
	if err != nil {
		return fmt.Errorf("failed to stop component: %w", err)
	}

	c.started = false
	return nil
}

// IsStarted checks if component is started
func (c *Component) IsStarted() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.started
}

// GetDependencies returns component dependencies
func (c *Component) GetDependencies() []string {
	return c.dependencies
}

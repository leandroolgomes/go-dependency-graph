package component

import (
	"fmt"
	"sort"
	"sync"
	"time"
)

// System manages all components and their lifecycle
type System struct {
	components map[string]*Component
	started    bool
	context    Context
	mu         sync.Mutex
}

// CreateSystem initializes a new system with the provided components
func CreateSystem(components map[string]*Component) *System {
	return &System{
		components: components,
		started:    false,
		context:    make(Context),
	}
}

// Start initializes all components in dependency order
func (s *System) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.started {
		return nil
	}

	systemStartTime := time.Now()

	// Check for cyclic dependencies
	if err := s.checkCyclicDependencies(); err != nil {
		return err
	}

	// Get components in order of dependencies
	orderedComponents, err := s.getOrderedComponents()
	if err != nil {
		return err
	}

	// Start components in order
	for _, name := range orderedComponents {
		component := s.components[name]
		
		// Create context with dependencies
		ctx := make(Context)
		for _, dep := range component.GetDependencies() {
			depComponent, exists := s.components[dep]
			if !exists {
				return fmt.Errorf("dependency %s not found for component %s", dep, name)
			}
			
			if !depComponent.IsStarted() {
				return fmt.Errorf("dependency %s not started for component %s", dep, name)
			}
			
			ctx[dep] = depComponent.instance
		}
		
		// Start the component
		lifecycle, err := component.Start(ctx)
		if err != nil {
			return fmt.Errorf("failed to start component %s: %w", name, err)
		}
		
		// Store the lifecycle instance in system context
		s.context[name] = lifecycle
	}
	
	systemElapsedTime := time.Since(systemStartTime)
	fmt.Printf("Total system initialization time: %v\n", systemElapsedTime)
	
	s.started = true
	return nil
}

// Stop gracefully shuts down all components in reverse dependency order
func (s *System) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.started {
		return nil
	}

	// Get components in order of dependencies
	orderedComponents, err := s.getOrderedComponents()
	if err != nil {
		return err
	}

	// Reverse the order for stopping
	for i, j := 0, len(orderedComponents)-1; i < j; i, j = i+1, j-1 {
		orderedComponents[i], orderedComponents[j] = orderedComponents[j], orderedComponents[i]
	}

	// Stop components in reverse order
	var lastErr error
	for _, name := range orderedComponents {
		component := s.components[name]
		if err := component.Stop(s.context); err != nil {
			lastErr = fmt.Errorf("failed to stop component %s: %w", name, err)
			// Continue stopping other components even if one fails
		}
	}

	s.started = false
	return lastErr
}

// GetContext returns the system context with all component results
func (s *System) GetContext() Context {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	// Create a copy to prevent external modification
	ctx := make(Context)
	for k, v := range s.context {
		ctx[k] = v
	}
	
	return ctx
}

// checkCyclicDependencies verifies that there are no cyclic dependencies
func (s *System) checkCyclicDependencies() error {
	visited := make(map[string]bool)
	recStack := make(map[string]bool)

	for name := range s.components {
		if !visited[name] {
			if s.isCyclic(name, visited, recStack) {
				return fmt.Errorf("cyclic dependency detected involving component %s", name)
			}
		}
	}

	return nil
}

// isCyclic is a helper function for cycle detection using DFS
func (s *System) isCyclic(name string, visited, recStack map[string]bool) bool {
	visited[name] = true
	recStack[name] = true

	component := s.components[name]
	for _, dep := range component.GetDependencies() {
		// Verificar se a dependência existe
		_, exists := s.components[dep]
		if !exists {
			// Dependência não encontrada, mas não é um ciclo
			continue
		}
		
		if !visited[dep] {
			if s.isCyclic(dep, visited, recStack) {
				return true
			}
		} else if recStack[dep] {
			return true
		}
	}

	recStack[name] = false
	return false
}

// getOrderedComponents returns components in dependency order
func (s *System) getOrderedComponents() ([]string, error) {
	// Build dependency graph
	graph := make(map[string][]string)
	inDegree := make(map[string]int)
	
	// Initialize all components with zero in-degree
	for name := range s.components {
		inDegree[name] = 0
		graph[name] = []string{}
	}
	
	// Calculate in-degree for each component
	for name, component := range s.components {
		for _, dep := range component.GetDependencies() {
			if _, exists := s.components[dep]; !exists {
				return nil, fmt.Errorf("dependency %s not found for component %s", dep, name)
			}
			graph[dep] = append(graph[dep], name)
			inDegree[name]++
		}
	}
	
	// Find all sources (nodes with in-degree 0)
	var queue []string
	for name, degree := range inDegree {
		if degree == 0 {
			queue = append(queue, name)
		}
	}
	
	// Topological sort
	var result []string
	for len(queue) > 0 {
		// Sort queue for deterministic order
		sort.Strings(queue)
		
		// Take first element
		current := queue[0]
		queue = queue[1:]
		result = append(result, current)
		
		// Reduce in-degree of neighbors
		for _, neighbor := range graph[current] {
			inDegree[neighbor]--
			if inDegree[neighbor] == 0 {
				queue = append(queue, neighbor)
			}
		}
	}
	
	// Check if all components were included
	if len(result) != len(s.components) {
		return nil, fmt.Errorf("cyclic dependency detected")
	}
	
	return result, nil
}

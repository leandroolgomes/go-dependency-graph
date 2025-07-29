package component

import (
	"errors"
	"testing"
)

// MockComponent implements the Lifecycle interface for testing
type MockComponent struct {
	Key         string
	StartCalled bool
	StopCalled  bool
	StartError  error
	StopError   error
}

func (m *MockComponent) Start(ctx Context) (Lifecycle, error) {
	m.StartCalled = true
	return m, m.StartError
}

func (m *MockComponent) Stop(ctx Context) error {
	m.StopCalled = true
	return m.StopError
}

func TestSystemStartStop(t *testing.T) {
	// Create mock components
	compA := &MockComponent{Key: "compA"}
	compB := &MockComponent{Key: "compB"}
	compC := &MockComponent{Key: "compC"}

	// Define components with dependencies
	components := map[string]*Component{
		"compA": Define("compA", compA),
		"compB": Define("compB", compB, "compA"),
		"compC": Define("compC", compC, "compA", "compB"),
	}

	// Create system
	system := CreateSystem(components)

	// Start the system
	if err := system.Start(); err != nil {
		t.Fatalf("Failed to start system: %v", err)
	}

	// Check that all components were started
	if !compA.StartCalled {
		t.Error("Component A was not started")
	}
	if !compB.StartCalled {
		t.Error("Component B was not started")
	}
	if !compC.StartCalled {
		t.Error("Component C was not started")
	}

	// Check context
	ctx := system.GetContext()
	if ctx["compA"] != compA {
		t.Errorf("Expected compA result to be the component itself")
	}
	if ctx["compB"] != compB {
		t.Errorf("Expected compB result to be the component itself")
	}
	if ctx["compC"] != compC {
		t.Errorf("Expected compC result to be the component itself")
	}

	// Stop the system
	if err := system.Stop(); err != nil {
		t.Fatalf("Failed to stop system: %v", err)
	}

	// Check that all components were stopped
	if !compA.StopCalled {
		t.Error("Component A was not stopped")
	}
	if !compB.StopCalled {
		t.Error("Component B was not stopped")
	}
	if !compC.StopCalled {
		t.Error("Component C was not stopped")
	}
}

func TestSystemStartError(t *testing.T) {
	// Create mock components with an error in compB
	compA := &MockComponent{Key: "compA"}
	compB := &MockComponent{Key: "compB", StartError: errors.New("start error")}

	// Define components with dependencies
	components := map[string]*Component{
		"compA": Define("compA", compA),
		"compB": Define("compB", compB, "compA"),
	}

	// Create system
	system := CreateSystem(components)

	// Start the system, should fail
	err := system.Start()
	if err == nil {
		t.Fatal("Expected system start to fail, but it succeeded")
	}

	// Check that compA was started but compB failed
	if !compA.StartCalled {
		t.Error("Component A was not started")
	}
	if !compB.StartCalled {
		t.Error("Component B was not started")
	}
}

func TestSystemCyclicDependency(t *testing.T) {
	// Create mock components with a cyclic dependency
	compA := &MockComponent{Key: "compA"}
	compB := &MockComponent{Key: "compB"}
	compC := &MockComponent{Key: "compC"}

	// Define components with a cyclic dependency: A -> B -> C -> A
	components := map[string]*Component{
		"compA": Define("compA", compA, "compC"),
		"compB": Define("compB", compB, "compA"),
		"compC": Define("compC", compC, "compB"),
	}

	// Create system
	system := CreateSystem(components)

	// Start the system, should fail due to cyclic dependency
	err := system.Start()
	if err == nil {
		t.Fatal("Expected system start to fail due to cyclic dependency, but it succeeded")
	}
}

func TestMissingDependency(t *testing.T) {
	// Create mock component with a missing dependency
	compA := &MockComponent{Key: "compA"}

	// Define component with a non-existent dependency
	components := map[string]*Component{
		"compA": Define("compA", compA, "nonExistent"),
	}

	// Create system
	system := CreateSystem(components)

	// Start the system, should fail due to missing dependency
	err := system.Start()
	if err == nil {
		t.Fatal("Expected system start to fail due to missing dependency, but it succeeded")
	}
}

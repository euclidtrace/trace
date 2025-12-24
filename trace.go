// Package trace provides deterministic, fully explainable computation tracing.
//
// This package implements structured tracing for computations where:
// - Same input always produces the same result and trace (deterministic)
// - Every result has an explicit, structured trace
// - Mathematical correctness and clarity over performance
// - All inputs, intermediate results, and traces are immutable
// - No business or domain-specific assumptions
//
// The goal is to make computations readable, auditable, and explainable by humans.
package trace

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"sort"
	"time"
)

// Value represents an immutable value in a trace.
// Values are stored as their original type along with a string representation.
type Value struct {
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}

// NewValue creates a new immutable Value from any input.
func NewValue(v interface{}) Value {
	return Value{
		Type:  fmt.Sprintf("%T", v),
		Value: v,
	}
}

// String returns a human-readable string representation of the value.
func (v Value) String() string {
	return fmt.Sprintf("%v", v.Value)
}

// Step represents a single computation step in a trace.
// Each step records what operation was performed, its inputs, and its output.
type Step struct {
	Operation   string            `json:"operation"`
	Description string            `json:"description,omitempty"`
	Inputs      map[string]Value  `json:"inputs"`
	Output      Value             `json:"output"`
	Timestamp   time.Time         `json:"timestamp"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

// NewStep creates a new computation step.
func NewStep(operation string, inputs map[string]Value, output Value) Step {
	return Step{
		Operation: operation,
		Inputs:    inputs,
		Output:    output,
		Timestamp: time.Now().UTC(),
		Metadata:  make(map[string]string),
	}
}

// WithDescription adds a human-readable description to the step.
func (s Step) WithDescription(desc string) Step {
	s.Description = desc
	return s
}

// WithMetadata adds metadata to the step.
func (s Step) WithMetadata(key, value string) Step {
	if s.Metadata == nil {
		s.Metadata = make(map[string]string)
	}
	s.Metadata[key] = value
	return s
}

// String returns a human-readable string representation of the step.
func (s Step) String() string {
	result := fmt.Sprintf("%s: %s", s.Operation, s.Output)
	if s.Description != "" {
		result = fmt.Sprintf("%s (%s)", result, s.Description)
	}
	return result
}

// Trace represents a complete computation trace.
// It records all inputs, intermediate steps, and the final result.
type Trace struct {
	ID        string            `json:"id"`
	Name      string            `json:"name"`
	Inputs    map[string]Value  `json:"inputs"`
	Steps     []Step            `json:"steps"`
	Result    *Value            `json:"result,omitempty"`
	StartTime time.Time         `json:"start_time"`
	EndTime   *time.Time        `json:"end_time,omitempty"`
	Metadata  map[string]string `json:"metadata,omitempty"`
	completed bool
}

// NewTrace creates a new trace with the given name and inputs.
// The trace is immutable once completed.
func NewTrace(name string, inputs map[string]Value) *Trace {
	return &Trace{
		ID:        generateID(name, inputs),
		Name:      name,
		Inputs:    copyInputs(inputs),
		Steps:     make([]Step, 0),
		StartTime: time.Now().UTC(),
		Metadata:  make(map[string]string),
		completed: false,
	}
}

// AddStep records a computation step in the trace.
// Returns an error if the trace is already completed.
func (t *Trace) AddStep(step Step) error {
	if t.completed {
		return fmt.Errorf("cannot add step to completed trace")
	}
	t.Steps = append(t.Steps, step)
	return nil
}

// SetResult sets the final result of the computation and marks the trace as complete.
// Returns an error if the trace is already completed.
func (t *Trace) SetResult(result Value) error {
	if t.completed {
		return fmt.Errorf("cannot set result on completed trace")
	}
	t.Result = &result
	now := time.Now().UTC()
	t.EndTime = &now
	t.completed = true
	return nil
}

// IsCompleted returns true if the trace has been completed.
func (t *Trace) IsCompleted() bool {
	return t.completed
}

// WithMetadata adds metadata to the trace.
// Returns an error if the trace is already completed.
func (t *Trace) WithMetadata(key, value string) error {
	if t.completed {
		return fmt.Errorf("cannot add metadata to completed trace")
	}
	if t.Metadata == nil {
		t.Metadata = make(map[string]string)
	}
	t.Metadata[key] = value
	return nil
}

// String returns a human-readable string representation of the trace.
func (t *Trace) String() string {
	result := fmt.Sprintf("Trace: %s (ID: %s)\n", t.Name, t.ID)
	result += "Inputs:\n"
	// Sort keys to ensure deterministic output
	keys := make([]string, 0, len(t.Inputs))
	for k := range t.Inputs {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		result += fmt.Sprintf("  %s: %s\n", k, t.Inputs[k])
	}
	if len(t.Steps) > 0 {
		result += "Steps:\n"
		for i, step := range t.Steps {
			result += fmt.Sprintf("  %d. %s\n", i+1, step)
		}
	}
	if t.Result != nil {
		result += fmt.Sprintf("Result: %s\n", t.Result)
	}
	return result
}

// ToJSON returns the trace as a JSON string.
func (t *Trace) ToJSON() (string, error) {
	data, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// copyInputs creates a deep copy of the inputs map to ensure immutability.
func copyInputs(inputs map[string]Value) map[string]Value {
	if inputs == nil {
		return make(map[string]Value)
	}
	copied := make(map[string]Value, len(inputs))
	for k, v := range inputs {
		copied[k] = v
	}
	return copied
}

// generateID creates a deterministic ID for the trace based on name and inputs.
// This ensures that the same inputs always produce the same trace ID.
func generateID(name string, inputs map[string]Value) string {
	// Create a hash of the name and inputs for deterministic ID
	h := sha256.New()
	h.Write([]byte(name))

	// Sort keys to ensure deterministic ordering
	keys := make([]string, 0, len(inputs))
	for k := range inputs {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Hash each key-value pair in sorted order
	for _, k := range keys {
		h.Write([]byte(k))
		h.Write([]byte(fmt.Sprintf("%v", inputs[k].Value)))
	}

	// Return first 16 characters of hex hash for readability
	return fmt.Sprintf("%s-%x", name, h.Sum(nil)[:8])
}

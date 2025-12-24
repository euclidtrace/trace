package trace

import (
	"encoding/json"
	"testing"
)

func TestNewValue(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		wantType string
	}{
		{"int", 42, "int"},
		{"float64", 3.14, "float64"},
		{"string", "hello", "string"},
		{"bool", true, "bool"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := NewValue(tt.input)
			if v.Type != tt.wantType {
				t.Errorf("NewValue().Type = %v, want %v", v.Type, tt.wantType)
			}
			if v.Value != tt.input {
				t.Errorf("NewValue().Value = %v, want %v", v.Value, tt.input)
			}
		})
	}
}

func TestValueString(t *testing.T) {
	tests := []struct {
		name  string
		value Value
		want  string
	}{
		{"int", NewValue(42), "42"},
		{"float", NewValue(3.14), "3.14"},
		{"string", NewValue("hello"), "hello"},
		{"bool", NewValue(true), "true"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.value.String(); got != tt.want {
				t.Errorf("Value.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewStep(t *testing.T) {
	inputs := map[string]Value{
		"a": NewValue(10),
		"b": NewValue(20),
	}
	output := NewValue(30)

	step := NewStep("add", inputs, output)

	if step.Operation != "add" {
		t.Errorf("NewStep().Operation = %v, want %v", step.Operation, "add")
	}
	if len(step.Inputs) != 2 {
		t.Errorf("NewStep().Inputs length = %v, want %v", len(step.Inputs), 2)
	}
	if step.Output.Value != 30 {
		t.Errorf("NewStep().Output.Value = %v, want %v", step.Output.Value, 30)
	}
	if step.Timestamp.IsZero() {
		t.Error("NewStep().Timestamp should not be zero")
	}
}

func TestStepWithDescription(t *testing.T) {
	step := NewStep("add", map[string]Value{}, NewValue(0))
	step = step.WithDescription("Addition operation")

	if step.Description != "Addition operation" {
		t.Errorf("Step.Description = %v, want %v", step.Description, "Addition operation")
	}
}

func TestStepWithMetadata(t *testing.T) {
	step := NewStep("add", map[string]Value{}, NewValue(0))
	step = step.WithMetadata("author", "test")

	if step.Metadata["author"] != "test" {
		t.Errorf("Step.Metadata['author'] = %v, want %v", step.Metadata["author"], "test")
	}
}

func TestNewTrace(t *testing.T) {
	inputs := map[string]Value{
		"x": NewValue(5),
		"y": NewValue(10),
	}

	tr := NewTrace("test-computation", inputs)

	if tr.Name != "test-computation" {
		t.Errorf("NewTrace().Name = %v, want %v", tr.Name, "test-computation")
	}
	if len(tr.Inputs) != 2 {
		t.Errorf("NewTrace().Inputs length = %v, want %v", len(tr.Inputs), 2)
	}
	if tr.ID == "" {
		t.Error("NewTrace().ID should not be empty")
	}
	if tr.StartTime.IsZero() {
		t.Error("NewTrace().StartTime should not be zero")
	}
	if tr.completed {
		t.Error("NewTrace() should not be completed initially")
	}
}

func TestTraceAddStep(t *testing.T) {
	tr := NewTrace("test", map[string]Value{})
	step := NewStep("operation", map[string]Value{}, NewValue(0))

	err := tr.AddStep(step)
	if err != nil {
		t.Errorf("AddStep() error = %v, want nil", err)
	}
	if len(tr.Steps) != 1 {
		t.Errorf("trace.Steps length = %v, want %v", len(tr.Steps), 1)
	}
}

func TestTraceAddStepAfterCompletion(t *testing.T) {
	tr := NewTrace("test", map[string]Value{})
	tr.SetResult(NewValue(0))

	step := NewStep("operation", map[string]Value{}, NewValue(0))
	err := tr.AddStep(step)

	if err == nil {
		t.Error("AddStep() after completion should return error")
	}
}

func TestTraceSetResult(t *testing.T) {
	tr := NewTrace("test", map[string]Value{})
	result := NewValue(42)

	err := tr.SetResult(result)
	if err != nil {
		t.Errorf("SetResult() error = %v, want nil", err)
	}
	if tr.Result == nil {
		t.Fatal("trace.Result should not be nil after SetResult()")
	}
	if tr.Result.Value != 42 {
		t.Errorf("trace.Result.Value = %v, want %v", tr.Result.Value, 42)
	}
	if !tr.completed {
		t.Error("trace should be completed after SetResult()")
	}
	if tr.EndTime == nil {
		t.Error("trace.EndTime should not be nil after SetResult()")
	}
}

func TestTraceSetResultTwice(t *testing.T) {
	tr := NewTrace("test", map[string]Value{})
	tr.SetResult(NewValue(1))

	err := tr.SetResult(NewValue(2))
	if err == nil {
		t.Error("SetResult() called twice should return error")
	}
}

func TestTraceIsCompleted(t *testing.T) {
	tr := NewTrace("test", map[string]Value{})

	if tr.IsCompleted() {
		t.Error("new trace should not be completed")
	}

	tr.SetResult(NewValue(0))

	if !tr.IsCompleted() {
		t.Error("trace should be completed after SetResult()")
	}
}

func TestTraceWithMetadata(t *testing.T) {
	tr := NewTrace("test", map[string]Value{})

	err := tr.WithMetadata("version", "1.0")
	if err != nil {
		t.Errorf("WithMetadata() error = %v, want nil", err)
	}
	if tr.Metadata["version"] != "1.0" {
		t.Errorf("trace.Metadata['version'] = %v, want %v", tr.Metadata["version"], "1.0")
	}
}

func TestTraceWithMetadataAfterCompletion(t *testing.T) {
	tr := NewTrace("test", map[string]Value{})
	tr.SetResult(NewValue(0))

	err := tr.WithMetadata("key", "value")
	if err == nil {
		t.Error("WithMetadata() after completion should return error")
	}
}

func TestTraceImmutability(t *testing.T) {
	// Test that modifying original inputs doesn't affect the trace
	inputs := map[string]Value{
		"x": NewValue(10),
	}

	tr := NewTrace("test", inputs)

	// Modify original inputs
	inputs["x"] = NewValue(20)
	inputs["y"] = NewValue(30)

	// Trace should still have original values
	if len(tr.Inputs) != 1 {
		t.Errorf("trace.Inputs length = %v, want %v", len(tr.Inputs), 1)
	}
	if tr.Inputs["x"].Value != 10 {
		t.Errorf("trace.Inputs['x'].Value = %v, want %v", tr.Inputs["x"].Value, 10)
	}
}

func TestTraceDeterministicID(t *testing.T) {
	inputs1 := map[string]Value{
		"a": NewValue(1),
		"b": NewValue(2),
	}
	inputs2 := map[string]Value{
		"a": NewValue(1),
		"b": NewValue(2),
	}

	tr1 := NewTrace("test", inputs1)
	tr2 := NewTrace("test", inputs2)

	if tr1.ID != tr2.ID {
		t.Errorf("same inputs should produce same ID: %v != %v", tr1.ID, tr2.ID)
	}
}

func TestTraceString(t *testing.T) {
	inputs := map[string]Value{
		"x": NewValue(5),
	}

	tr := NewTrace("test-trace", inputs)
	tr.AddStep(NewStep("step1", map[string]Value{"input": NewValue(5)}, NewValue(10)))
	tr.SetResult(NewValue(10))

	str := tr.String()
	if str == "" {
		t.Error("Trace.String() should not be empty")
	}
	// Basic check that it contains expected elements
	if len(str) < 10 {
		t.Error("Trace.String() output seems too short")
	}
}

func TestTraceToJSON(t *testing.T) {
	inputs := map[string]Value{
		"x": NewValue(5),
	}

	tr := NewTrace("test-trace", inputs)
	tr.AddStep(NewStep("step1", map[string]Value{"input": NewValue(5)}, NewValue(10)))
	tr.SetResult(NewValue(10))

	jsonStr, err := tr.ToJSON()
	if err != nil {
		t.Errorf("ToJSON() error = %v, want nil", err)
	}
	if jsonStr == "" {
		t.Error("ToJSON() should not return empty string")
	}

	// Verify it's valid JSON
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		t.Errorf("ToJSON() did not produce valid JSON: %v", err)
	}
}

func TestCompleteComputationTrace(t *testing.T) {
	// Simulate a complete computation: (a + b) * c
	inputs := map[string]Value{
		"a": NewValue(10),
		"b": NewValue(20),
		"c": NewValue(3),
	}

	tr := NewTrace("multiply-sum", inputs)

	// Step 1: a + b
	step1Inputs := map[string]Value{
		"a": NewValue(10),
		"b": NewValue(20),
	}
	step1 := NewStep("add", step1Inputs, NewValue(30))
	step1 = step1.WithDescription("Add a and b")
	tr.AddStep(step1)

	// Step 2: result * c
	step2Inputs := map[string]Value{
		"sum": NewValue(30),
		"c":   NewValue(3),
	}
	step2 := NewStep("multiply", step2Inputs, NewValue(90))
	step2 = step2.WithDescription("Multiply sum by c")
	tr.AddStep(step2)

	// Set final result
	tr.SetResult(NewValue(90))

	// Verify trace structure
	if !tr.IsCompleted() {
		t.Error("trace should be completed")
	}
	if len(tr.Steps) != 2 {
		t.Errorf("trace should have 2 steps, got %v", len(tr.Steps))
	}
	if tr.Result.Value != 90 {
		t.Errorf("trace.Result.Value = %v, want %v", tr.Result.Value, 90)
	}

	// Verify JSON can be generated
	_, err := tr.ToJSON()
	if err != nil {
		t.Errorf("complete trace ToJSON() error = %v", err)
	}

	// Verify string representation
	str := tr.String()
	if str == "" {
		t.Error("complete trace String() should not be empty")
	}
}

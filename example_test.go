package trace_test

import (
	"fmt"

	"github.com/EuclidOLAP/trace"
)

// Example demonstrates basic usage of the trace package.
func Example() {
	// Define inputs for a computation
	inputs := map[string]trace.Value{
		"a": trace.NewValue(10),
		"b": trace.NewValue(20),
	}

	// Create a new trace
	tr := trace.NewTrace("addition", inputs)

	// Record a computation step
	step := trace.NewStep(
		"add",
		map[string]trace.Value{
			"a": trace.NewValue(10),
			"b": trace.NewValue(20),
		},
		trace.NewValue(30),
	)
	tr.AddStep(step)

	// Set the final result
	tr.SetResult(trace.NewValue(30))

	// Output the trace
	fmt.Println(tr.String())
}

// ExampleTrace_complexComputation demonstrates a multi-step computation trace.
func ExampleTrace_complexComputation() {
	// Computation: ((a + b) * c) - d
	inputs := map[string]trace.Value{
		"a": trace.NewValue(5),
		"b": trace.NewValue(3),
		"c": trace.NewValue(4),
		"d": trace.NewValue(2),
	}

	tr := trace.NewTrace("complex-computation", inputs)

	// Step 1: Add a and b
	step1 := trace.NewStep(
		"add",
		map[string]trace.Value{
			"a": trace.NewValue(5),
			"b": trace.NewValue(3),
		},
		trace.NewValue(8),
	).WithDescription("Add a and b")
	tr.AddStep(step1)

	// Step 2: Multiply result by c
	step2 := trace.NewStep(
		"multiply",
		map[string]trace.Value{
			"sum": trace.NewValue(8),
			"c":   trace.NewValue(4),
		},
		trace.NewValue(32),
	).WithDescription("Multiply sum by c")
	tr.AddStep(step2)

	// Step 3: Subtract d
	step3 := trace.NewStep(
		"subtract",
		map[string]trace.Value{
			"product": trace.NewValue(32),
			"d":       trace.NewValue(2),
		},
		trace.NewValue(30),
	).WithDescription("Subtract d from product")
	tr.AddStep(step3)

	// Set final result
	tr.SetResult(trace.NewValue(30))

	fmt.Printf("Computation completed: %d\n", tr.Result.Value)
	fmt.Printf("Total steps: %d\n", len(tr.Steps))
	// Output:
	// Computation completed: 30
	// Total steps: 3
}

// ExampleTrace_withMetadata demonstrates adding metadata to traces.
func ExampleTrace_withMetadata() {
	inputs := map[string]trace.Value{
		"value": trace.NewValue(100),
	}

	tr := trace.NewTrace("percentage-calculation", inputs)
	tr.WithMetadata("unit", "percent")
	tr.WithMetadata("precision", "2")

	step := trace.NewStep(
		"divide",
		map[string]trace.Value{
			"numerator":   trace.NewValue(100),
			"denominator": trace.NewValue(100),
		},
		trace.NewValue(1.0),
	).WithMetadata("explanation", "Convert to decimal")

	tr.AddStep(step)
	tr.SetResult(trace.NewValue(1.0))

	fmt.Println("Trace completed with metadata")
	// Output: Trace completed with metadata
}

// ExampleTrace_toJSON demonstrates JSON serialization of traces.
func ExampleTrace_toJSON() {
	inputs := map[string]trace.Value{
		"x": trace.NewValue(42),
	}

	tr := trace.NewTrace("identity", inputs)
	tr.SetResult(trace.NewValue(42))

	jsonStr, err := tr.ToJSON()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Println("JSON generated successfully")
	fmt.Printf("JSON length > 250 bytes: %v\n", len(jsonStr) > 250)
	// Output:
	// JSON generated successfully
	// JSON length > 250 bytes: true
}

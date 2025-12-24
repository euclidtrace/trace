# trace

Deterministic, fully explainable computation tracing for Go.

Part of [Euclid Analysis](https://github.com/EuclidOLAP), the `trace` package provides structured tracing for computations where every result must be explainable and reproducible.

## Purpose

This library is designed for systems that require:
- **Determinism**: Same input always produces the same result and trace
- **Full Explainability**: Every result has an explicit, structured trace
- **Mathematical Correctness**: Clarity and correctness over performance
- **Immutability**: Inputs, intermediate results, and traces cannot be modified
- **Domain Independence**: No business or domain-specific assumptions

The goal is to make computations readable, auditable, and explainable by humans.

## Principles

1. **Deterministic**: Same input → same result and trace, always
2. **Explicit**: Every computation step is recorded with inputs and outputs
3. **Immutable**: Once created, values and completed traces cannot be changed
4. **Structured**: Traces are not logs; they are structured data
5. **Human-Readable**: All traces can be viewed as text or JSON
6. **Domain-Agnostic**: No assumptions about what is being computed

## Installation

```bash
go get github.com/EuclidOLAP/trace
```

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/EuclidOLAP/trace"
)

func main() {
    // Define inputs
    inputs := map[string]trace.Value{
        "a": trace.NewValue(10),
        "b": trace.NewValue(20),
    }

    // Create a trace
    tr := trace.NewTrace("addition", inputs)

    // Record computation step
    step := trace.NewStep(
        "add",
        map[string]trace.Value{
            "a": trace.NewValue(10),
            "b": trace.NewValue(20),
        },
        trace.NewValue(30),
    )
    tr.AddStep(step)

    // Set final result
    tr.SetResult(trace.NewValue(30))

    // View the trace
    fmt.Println(tr.String())
}
```

## Core Concepts

### Value

An immutable value with type information:

```go
v := trace.NewValue(42)
fmt.Println(v.String())  // "42"
fmt.Println(v.Type)      // "int"
```

### Step

A single computation step recording an operation, its inputs, and output:

```go
step := trace.NewStep(
    "multiply",
    map[string]trace.Value{
        "x": trace.NewValue(5),
        "y": trace.NewValue(3),
    },
    trace.NewValue(15),
)

// Add description and metadata
step = step.WithDescription("Multiply x by y")
step = step.WithMetadata("precision", "exact")
```

### Trace

A complete computation trace with all inputs, steps, and result:

```go
tr := trace.NewTrace("computation", inputs)
tr.AddStep(step1)
tr.AddStep(step2)
tr.SetResult(finalValue)

// Once completed, the trace is immutable
fmt.Println(tr.IsCompleted())  // true
```

## Examples

### Multi-Step Computation

```go
// Compute: ((a + b) * c) - d
inputs := map[string]trace.Value{
    "a": trace.NewValue(5),
    "b": trace.NewValue(3),
    "c": trace.NewValue(4),
    "d": trace.NewValue(2),
}

tr := trace.NewTrace("complex-computation", inputs)

// Step 1: Add
step1 := trace.NewStep("add",
    map[string]trace.Value{
        "a": trace.NewValue(5),
        "b": trace.NewValue(3),
    },
    trace.NewValue(8),
).WithDescription("Add a and b")
tr.AddStep(step1)

// Step 2: Multiply
step2 := trace.NewStep("multiply",
    map[string]trace.Value{
        "sum": trace.NewValue(8),
        "c":   trace.NewValue(4),
    },
    trace.NewValue(32),
).WithDescription("Multiply sum by c")
tr.AddStep(step2)

// Step 3: Subtract
step3 := trace.NewStep("subtract",
    map[string]trace.Value{
        "product": trace.NewValue(32),
        "d":       trace.NewValue(2),
    },
    trace.NewValue(30),
).WithDescription("Subtract d from product")
tr.AddStep(step3)

tr.SetResult(trace.NewValue(30))
```

### JSON Export

```go
tr := trace.NewTrace("calculation", inputs)
// ... add steps ...
tr.SetResult(result)

jsonStr, err := tr.ToJSON()
if err != nil {
    log.Fatal(err)
}
fmt.Println(jsonStr)
```

### Metadata

```go
tr := trace.NewTrace("analysis", inputs)
tr.WithMetadata("version", "1.0")
tr.WithMetadata("author", "system")

step := trace.NewStep("normalize", inputs, output)
step = step.WithMetadata("method", "min-max")
tr.AddStep(step)
```

## Design Decisions

### Why Immutable?

Immutability ensures that traces are reproducible and trustworthy. Once a computation is complete, its trace cannot be altered, making it suitable for audit trails and verification.

### Why Structured Traces, Not Logs?

Logs are unstructured text for debugging. Traces are structured data that explain *what* was computed, *how* it was computed, and *why* the result is what it is.

### Why Deterministic IDs?

Traces generate deterministic IDs based on their name and inputs. This allows the same computation to be identified across runs, making it easier to cache results or verify consistency.

### Why No Performance Focus?

This library prioritizes correctness and clarity. If you need high-performance tracing, consider other tools. This library is for computations where being right and explainable matters more than being fast.

## Testing

Run tests:

```bash
go test ./...
```

Run tests with coverage:

```bash
go test -cover ./...
```

## Use Cases

- Mathematical computations that require audit trails
- Financial calculations that must be explainable
- Scientific computations that need reproducibility
- Regulatory compliance where every decision must be justified
- Educational tools that explain step-by-step computations
- Debugging complex algorithms by examining their traces

## Not For

- High-frequency, performance-critical systems
- Logging application events (use a logging library)
- Distributed tracing (use OpenTelemetry or similar)
- Business logic with domain-specific requirements

## License

MIT

## Contributing

This library follows strict principles. Contributions are welcome if they:
- Maintain determinism
- Preserve immutability
- Keep the API simple and clear
- Add no domain-specific features
- Include tests and documentation
# Technical Assignment - Hemant

## Overview

A **task graph** is a directed, potentially cyclic graph of nodes that orchestrates workflow execution. Every task graph must start at a `start` node and terminate at an `end` node. The graph supports complex routing patterns where:

- Multiple edges can flow into a single node
- Multiple edges can flow out of a single node
- Nodes can dynamically decide at runtime which subset of their outgoing edges to execute

### Example Routing Patterns

**Pattern 1: Conditional Branching**

```
start → a → b → [c, d, e] → f → end

```

Here, node `b` can decide to route to any subset of nodes `c`, `d`, or `e` based on runtime conditions.

**Pattern 2: Cyclic Processing**

```
start → a → b → c → [b, d] → d → end

```

Here, node `c` can decide whether to loop back to `b` for further analysis or continue forward to `d`.

### System Context

In the actual system, nodes have particular semantics such as:

- **Agent invocations**: AI agents that process requests and make decisions
- **Deterministic tool calling**: Functions that perform specific operations

Your work will focus on building the **execution runtime** for task graphs. This includes:

- **Execution semantics**: Managing the flow and state of the task graph
- **Resource scheduling**: Ensuring appropriate computational resources (e.g., code interpreter pods) are available when nodes need them, while finding an efficient compromise between latency, throughput, and resource utilization
- **Telemetry and observability**: Comprehensive logging, metrics, and tracing for offline debugging of individual task graph runs and pattern analysis across multiple executions
- **Authentication and authorization**: RBAC (Role-Based Access Control) for request scoping
- **Durable execution**: Handling process crashes, partial failures, and recovery scenarios to ensure workflows can resume from interruption points
- **Population-scale execution**: Managing thousands of concurrent task graph executions with proper coordination and orchestration

## Assignment Tasks

### Task 1: Graph Invariants

Identify and document the invariants that a valid task graph must satisfy. These should be properties that can be verified statically (before execution). Consider questions like:

- What connectivity requirements must be met?
- Are there constraints on node relationships?
- How should cycles be handled?
- What constitutes an "orphaned" node?

**Example invariant**: All paths from `start` must eventually reach `end` (no orphan nodes).

### Task 2: Validation Implementation

Write a function to validate that a given graph representation satisfies your identified invariants.

**Requirements**:

- Use any reasonable programming language (Python, Java, Go, Rust, etc. - not PHP or Malboge or anything else that requires divine intervention to parse or debug)
- You may choose your input representation (we use adjacency lists, but you're free to pick any reasonable format)
- Include comprehensive test cases covering both valid and invalid graphs
- You are encouraged to use LLMs and any available standard libraries to implement your solution

**Function signature example** (adapt to your chosen language):

```go
// CheckInvariantViolations validates a task graph against defined invariants.
// Returns a list of invariant violations. Empty list means valid graph.
func CheckInvariantViolations(graph map[string][]string) []string {
    // Implementation here
    return nil
}

// Alternative name suggestion:
// func ValidateTaskGraph(graph map[string][]string) []string
```

**Test cases**: Provide comprehensive test inputs that validate your invariant checking. These can be in the form of unit tests, sample JSON files, or any other appropriate format of your choosing. Your test cases should cover both valid and invalid graphs, including edge cases.

### Task 3: Durable Execution Design

Write a short design document (a couple of paragraphs) suggesting potential durable execution semantics for task graphs.

**Key questions to address**:

- How should the system handle process crashes during execution?
- What state needs to be persisted for recovery?
- How can partial execution be resumed?
- What are the trade-offs between different recovery approaches?

**Consider approaches such as**:

- Checkpoint-based recovery
- Event sourcing
- Workflow engines (Temporal, Cadence, etc.)
- Database-backed state management
- Idempotent execution patterns

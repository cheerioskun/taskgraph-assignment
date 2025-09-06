# TaskGraph Validator and Visualizer

A Go-based task graph validation system with p5.js visualization, implementing the requirements from the technical assignment.

## Features

- **Graph Validation**: Validates task graphs against defined invariants
- **Interactive Visualization**: p5.js-powered web interface for exploring graphs
- **CLI Interface**: Command-line tools for validation and testing
- **Sample Graphs**: Predefined examples including valid and invalid cases
- **Real-time Feedback**: Immediate validation results with detailed error messages

## Quick Start

1. **Initialize the Go module**:
   ```bash
   go mod tidy
   ```

2. **Start the visualization server**:
   ```bash
   go run . server
   ```

3. **Open your browser** to `http://localhost:8080`

## CLI Usage

```bash
# Start web server for visualization
go run . server

# Validate a graph from JSON file
go run . validate sample.json

# Show sample graphs with validation status
go run . samples

# Run basic validation tests
go run . test

# Run comprehensive tests
go test
```

## Graph Invariants

The system validates the following invariants:

1. **Single Start/End Nodes**: Exactly one `start` and one `end` node required
2. **Connectivity**: All nodes must be reachable from `start` and must be able to reach `end`
3. **No Isolated Components**: Every node must be in R(start) ∩ RR(end)
4. **Unique Edges**: No duplicate edges between the same pair of nodes
5. **No Orphaned Nodes**: All nodes must be part of valid execution paths

## Graph JSON Format

```json
{
  "nodes": [
    {
      "id": "start",
      "type": "start",
      "metadata": {"x": 100, "y": 150}
    },
    {
      "id": "process", 
      "type": "agent",
      "metadata": {"x": 300, "y": 150}
    },
    {
      "id": "end",
      "type": "end", 
      "metadata": {"x": 500, "y": 150}
    }
  ],
  "edges": [
    {"from": "start", "to": "process"},
    {"from": "process", "to": "end"}
  ]
}
```

### Node Types
- `start`: Entry point of the graph
- `end`: Exit point of the graph  
- `agent`: AI agent invocation
- `tool`: Deterministic tool execution

### Optional Metadata
- `x`, `y`: Position coordinates for visualization
- `condition`: Edge condition for conditional routing

## Visualization Features

- **Interactive Nodes**: Drag nodes to reposition them
- **Pan & Zoom**: Navigate large graphs easily
- **Color Coding**: Different colors for each node type
- **Error Highlighting**: Invalid nodes highlighted in purple
- **Real-time Validation**: Validate graphs with immediate feedback
- **Sample Graphs**: Pre-loaded examples including edge cases

## API Endpoints

- `GET /api/sample-graphs`: Returns available sample graphs
- `POST /api/validate`: Validates a submitted graph

## Architecture

```
├── types.go          # Core data structures
├── validator.go      # Graph validation logic
├── validator_test.go # Comprehensive test suite
├── server.go         # HTTP server and API handlers
├── main.go           # CLI interface and entry point
└── web/
    ├── index.html    # Visualization interface
    └── app.js        # p5.js visualization logic
```

## Durable Execution Design

The system is designed with future durable execution in mind:

### Checkpoint-based Recovery
- **State Persistence**: Execution state includes frontier tracking, completed nodes, and checkpoints
- **Recovery Points**: Each node completion creates a recovery checkpoint
- **Partial Resume**: Failed executions can resume from the last successful checkpoint

### Execution Model
- **Frontier Tracking**: Active execution frontier maintained as a set of ready-to-execute nodes
- **Dynamic Routing**: Nodes can decide at runtime which outgoing edges to activate
- **Concurrent Execution**: Multiple paths can execute simultaneously where dependencies allow

### Trade-offs Considered
- **Latency vs Durability**: Checkpoint frequency affects both recovery granularity and performance
- **Storage vs Memory**: State can be persisted to database or maintained in-memory for speed
- **Consistency vs Availability**: Strong consistency ensures correctness but may impact throughput

## Sample Graphs Included

1. **Simple Linear**: Basic start→process→end flow
2. **Branching**: Conditional routing with multiple paths
3. **Cyclic Processing**: Feedback loops with conditional continuation
4. **Invalid Examples**: Demonstrates various invariant violations

## Testing

Run the comprehensive test suite:
```bash
go test -v
```

Tests cover:
- Valid graph patterns (linear, branching, cyclic)
- All invariant violation types
- Edge cases and boundary conditions
- Graph traversal algorithms

## Future Extensions

The scaffold supports extension for:
- Runtime execution engine
- Resource scheduling
- Authentication/authorization
- Distributed execution
- Telemetry and observability
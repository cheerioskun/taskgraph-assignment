# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a TaskGraph validator and visualizer for Sarvam AI. It implements a directed graph with specific invariants that represent task execution workflows. The system includes validation logic and a web API for graph validation.

## Core Architecture

### Graph Structure
- **TaskGraph**: Main data structure representing a directed graph with nodes and edges
- **Node Types**: `start`, `end`, `agent`, `tool` - each graph must have exactly one start and end node
- **Validation System**: Pluggable validation functions that check graph invariants
- **Adjacency Lists**: Both forward and reverse adjacency lists are maintained for efficient traversal

### Key Components
- `types.go`: Core data structures (TaskGraph, TaskNode, TaskEdge, ValidationError)
- `graph.go`: Graph construction and validation orchestration 
- `validations.go`: Individual validation functions for graph invariants
- `server.go`: HTTP API server with `/validate` endpoint
- `errors.go`: Predefined error types for validation failures

### Validation Rules
The system validates these invariants:
1. Exactly one start node exists
2. Exactly one end node exists  
3. No duplicate edges
4. Connectivity: all nodes reachable from start and can reach end

Validation errors are categorized as:
- **Isolated nodes**: Neither reachable from start nor can reach end
- **Unrunnable nodes**: Cannot be reached from start
- **Orphaned nodes**: Cannot reach the end (warning level)

## Commands

### Build and Run
```bash
go run .                    # Start HTTP server on port 8080
go build                    # Build binary
```

### Testing
```bash
go test                     # Run all tests
go test -v                  # Run tests with verbose output
go test -run TestValidateConnectivity  # Run specific test
```

### Development
```bash
go mod tidy                 # Clean up dependencies
go fmt ./...                # Format code
```

## API Usage

The server exposes a POST `/validate` endpoint that accepts:
```json
{
  "nodes": [{"id": "start", "type": "start"}, ...],
  "edges": [{"from": "start", "to": "end"}, ...]
}
```

Returns validation results with error details including affected nodes/edges.
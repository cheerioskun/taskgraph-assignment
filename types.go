package main

// NodeType represents the type of a task node
type NodeType string

const (
	StartNode NodeType = "start"
	EndNode   NodeType = "end"
	AgentNode NodeType = "agent"
	ToolNode  NodeType = "tool"
)

type TaskNode struct {
	ID       string         `json:"id"`
	Type     NodeType       `json:"type"`
	Metadata map[string]any `json:"metadata,omitempty"`
}

type TaskEdge struct {
	From string `json:"from"`
	To   string `json:"to"`
}

type TaskGraph struct {
	StartNodeID          string              `json:"start_node_id,omitempty"`
	EndNodeID            string              `json:"end_node_id,omitempty"`
	Nodes                []TaskNode          `json:"nodes"`
	Edges                []TaskEdge          `json:"edges"`
	AdjacencyList        map[string][]string `json:"-"`
	ReverseAdjacencyList map[string][]string `json:"-"`
}

type ValidationFunction func(graph *TaskGraph) []ValidationError

// ValidationError represents a graph validation error
type ValidationError struct {
	error
	Nodes []string   `json:"nodes,omitempty"`
	Edges []TaskEdge `json:"edges,omitempty"`
}

// ExecutionState represents the runtime state of graph execution
type ExecutionState struct {
	Frontier  []string `json:"frontier"`
	Completed []string `json:"completed"`
	Active    []string `json:"active"`
}

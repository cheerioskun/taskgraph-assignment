package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateConnectivity(t *testing.T) {
	tests := []struct {
		name           string
		nodes          []TaskNode
		edges          []TaskEdge
		expectedErrors []error
	}{
		{
			name: "valid connected graph",
			nodes: []TaskNode{
				{ID: "start", Type: StartNode},
				{ID: "middle", Type: AgentNode},
				{ID: "end", Type: EndNode},
			},
			edges: []TaskEdge{
				{From: "start", To: "middle"},
				{From: "middle", To: "end"},
			},
			expectedErrors: []error{},
		},
		{
			name: "graph with cycle",
			nodes: []TaskNode{
				{ID: "start", Type: StartNode},
				{ID: "A", Type: AgentNode},
				{ID: "B", Type: ToolNode},
				{ID: "end", Type: EndNode},
			},
			edges: []TaskEdge{
				{From: "start", To: "A"},
				{From: "A", To: "B"},
				{From: "B", To: "A"}, // Cycle here
				{From: "B", To: "end"},
			},
			expectedErrors: []error{},
		},
		{
			name: "graph with isolated node",
			nodes: []TaskNode{
				{ID: "start", Type: StartNode},
				{ID: "isolated", Type: AgentNode},
				{ID: "end", Type: EndNode},
			},
			edges: []TaskEdge{
				{From: "start", To: "end"},
			},
			expectedErrors: []error{ErrIsolatedNodes},
		},
		{
			name: "graph with unrunnable node",
			nodes: []TaskNode{
				{ID: "start", Type: StartNode},
				{ID: "unrunnable", Type: AgentNode},
				{ID: "end", Type: EndNode},
			},
			edges: []TaskEdge{
				{From: "start", To: "end"},
				{From: "unrunnable", To: "end"},
			},
			expectedErrors: []error{ErrUnrunnableNodes},
		},
		{
			name: "graph with orphaned node",
			nodes: []TaskNode{
				{ID: "start", Type: StartNode},
				{ID: "orphaned", Type: AgentNode},
				{ID: "end", Type: EndNode},
			},
			edges: []TaskEdge{
				{From: "start", To: "orphaned"},
				{From: "start", To: "end"},
			},
			expectedErrors: []error{ErrOrphanedNodes},
		},
		{
			name: "graph with cycle and multiple connectivity issues",
			nodes: []TaskNode{
				{ID: "start", Type: StartNode},
				{ID: "isolated1", Type: AgentNode},
				{ID: "isolated2", Type: ToolNode},
				{ID: "A", Type: AgentNode},
				{ID: "orphaned", Type: AgentNode},
				{ID: "unrunnable", Type: ToolNode},
				{ID: "end", Type: EndNode},
			},
			edges: []TaskEdge{
				{From: "start", To: "orphaned"},
				{From: "start", To: "A"},
				{From: "A", To: "end"},
				{From: "A", To: "A"}, // Cycle
				{From: "isolated1", To: "isolated2"},
				{From: "unrunnable", To: "end"},
			},
			expectedErrors: []error{ErrIsolatedNodes, ErrUnrunnableNodes, ErrOrphanedNodes},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			graph := NewTaskGraph(tt.nodes, tt.edges)
			errors := ValidateConnectivity(graph)

			assert.Len(t, errors, len(tt.expectedErrors))

			for i, expectedErr := range tt.expectedErrors {
				assert.Equal(t, expectedErr, errors[i].error)
			}
		})
	}
}

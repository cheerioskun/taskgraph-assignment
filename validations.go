package main

import (
	"github.com/hashicorp/go-set/v3"
)

func ValidateStartNodeExists(graph *TaskGraph) []ValidationError {
	var errors []ValidationError
	startNodeIds := []string{}
	for _, node := range graph.Nodes {
		if node.Type == StartNode {
			startNodeIds = append(startNodeIds, node.ID)
		}
	}
	if len(startNodeIds) == 0 {
		errors = append(errors, ValidationError{
			error: ErrMissingStart,
		})
	} else if len(startNodeIds) > 1 {
		errors = append(errors, ValidationError{
			error: ErrMultipleStart,
			Nodes: startNodeIds,
		})
	}
	return errors
}

func ValidateEndNodeExists(graph *TaskGraph) []ValidationError {
	var errors []ValidationError
	endNodeIds := []string{}
	for _, node := range graph.Nodes {
		if node.Type == EndNode {
			endNodeIds = append(endNodeIds, node.ID)
		}
	}
	if len(endNodeIds) == 0 {
		errors = append(errors, ValidationError{
			error: ErrMissingEnd,
		})
	} else if len(endNodeIds) > 1 {
		errors = append(errors, ValidationError{
			error: ErrMultipleEnd,
			Nodes: endNodeIds,
		})
	}
	return errors
}

func ValidateNoDuplicateEdges(graph *TaskGraph) []ValidationError {
	var errors []ValidationError
	seen := make(map[TaskEdge]bool)
	for _, edge := range graph.Edges {
		if seen[edge] {
			errors = append(errors, ValidationError{
				error: ErrDuplicateEdge,
				Edges: []TaskEdge{edge},
			})
		} else {
			seen[edge] = true
		}
	}
	return errors
}

// We create three types of unconnectivity distinctions
// - Unreachable from start (nodes that cannot be reached from the start node): Unrunnable
// - Cannot reach end (nodes that cannot reach the end node): Orphans
// - Not part of component (nodes that cannot be reached from start and cannot reach end either): Isolated
func ValidateConnectivity(graph *TaskGraph) []ValidationError {
	// var errors []ValidationError

	start := graph.StartNodeID
	// Create reachability set
	queue := []string{start}
	reachableFromStart := set.New[string](len(graph.Nodes))
	reachableFromStart.Insert(start)
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		reachableFromStart.Insert(current)
		for _, neighbor := range graph.AdjacencyList[current] {
			if !reachableFromStart.Contains(neighbor) {
				queue = append(queue, neighbor)
			}
		}
	}
	end := graph.EndNodeID
	// Create reverse reachability set
	queue = []string{end}
	reachableToEnd := set.New[string](len(graph.Nodes))
	reachableToEnd.Insert(end)
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		reachableToEnd.Insert(current)
		for _, neighbor := range graph.ReverseAdjacencyList[current] {
			if !reachableToEnd.Contains(neighbor) {
				queue = append(queue, neighbor)
			}
		}
	}

	totalSet := set.New[string](len(graph.Nodes))
	for _, node := range graph.Nodes {
		totalSet.Insert(node.ID)
	}
	// Remove isolated nodes first. These are nodes that are neither reachable from start nor can reach end.
	var errors []ValidationError
	isolated := totalSet.Difference(reachableFromStart.Union(reachableToEnd))
	if isolated.Size() > 0 {
		errors = append(errors, ValidationError{
			error: ErrIsolatedNodes,
			Nodes: isolated.Slice(),
		})
	}
	totalViableSet := totalSet.Difference(isolated)
	unrunnable := totalViableSet.Difference(reachableFromStart)
	if unrunnable.Size() > 0 {
		errors = append(errors, ValidationError{
			error: ErrUnrunnableNodes,
			Nodes: unrunnable.Slice(),
		})
	}
	orphaned := totalViableSet.Difference(reachableToEnd)
	if orphaned.Size() > 0 {
		errors = append(errors, ValidationError{
			error: ErrOrphanedNodes,
			Nodes: orphaned.Slice(),
		})
	}
	return errors
}

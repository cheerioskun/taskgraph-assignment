package main

var (
	validations = []ValidationFunction{
		ValidateStartNodeExists,
		ValidateEndNodeExists,
		ValidateNoDuplicateEdges,
		ValidateConnectivity,
	}
)

func NewTaskGraph(nodes []TaskNode, edges []TaskEdge) *TaskGraph {

	var startNodeID string
	var endNodeID string
	// Identify start and end nodes
	for _, node := range nodes {
		if node.Type == StartNode {
			startNodeID = node.ID
		} else if node.Type == EndNode {
			endNodeID = node.ID
		}
	}
	tg := &TaskGraph{
		StartNodeID: startNodeID,
		EndNodeID:   endNodeID,
		Nodes:       nodes,
		Edges:       edges,
	}
	tg.buildAdjacencyLists()
	return tg
}

func (tg *TaskGraph) buildAdjacencyLists() {
	tg.AdjacencyList = make(map[string][]string)
	tg.ReverseAdjacencyList = make(map[string][]string)
	for _, edge := range tg.Edges {
		tg.AdjacencyList[edge.From] = append(tg.AdjacencyList[edge.From], edge.To)
		tg.ReverseAdjacencyList[edge.To] = append(tg.ReverseAdjacencyList[edge.To], edge.From)
	}
}

// CheckInvariantViolations validates a task graph against defined invariants
func CheckInvariantViolations(graph *TaskGraph) []ValidationError {
	for _, validate := range validations {
		errs := validate(graph)
		if len(errs) > 0 {
			return errs
		}
	}
	return nil
}

// IsValidTaskGraph returns true if the graph has no invariant violations
func IsValidTaskGraph(graph *TaskGraph) bool {
	return len(CheckInvariantViolations(graph)) == 0
}

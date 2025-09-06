package main

import (
	"encoding/json"
	"net/http"
)

type ValidationRequest struct {
	Nodes []TaskNode `json:"nodes"`
	Edges []TaskEdge `json:"edges"`
}

type ValidationResponse struct {
	Valid  bool       `json:"valid"`
	Errors []APIError `json:"errors"`
}

type APIError struct {
	Type     string     `json:"type"`
	Message  string     `json:"message"`
	Nodes    []string   `json:"nodes,omitempty"`
	Edges    []TaskEdge `json:"edges,omitempty"`
	Severity string     `json:"severity"`
}

func validateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ValidationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	graph := NewTaskGraph(req.Nodes, req.Edges)
	validationErrors := CheckInvariantViolations(graph)

	apiErrors := make([]APIError, 0, len(validationErrors))
	for _, err := range validationErrors {
		apiErrors = append(apiErrors, convertToAPIError(err))
	}

	response := ValidationResponse{
		Valid:  len(validationErrors) == 0,
		Errors: apiErrors,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func convertToAPIError(err ValidationError) APIError {
	severity := "error"
	message := getErrorMessage(err.error)

	if err.error == ErrOrphanedNodes {
		severity = "warning"
	}

	return APIError{
		Type:     err.Error(),
		Message:  message,
		Nodes:    err.Nodes,
		Edges:    err.Edges,
		Severity: severity,
	}
}

func getErrorMessage(err error) string {
	switch err {
	case ErrMissingStart:
		return "Graph is missing a start node"
	case ErrMissingEnd:
		return "Graph is missing an end node"
	case ErrMultipleStart:
		return "Graph has multiple start nodes"
	case ErrMultipleEnd:
		return "Graph has multiple end nodes"
	case ErrUnrunnableNodes:
		return "Nodes cannot be reached from the start"
	case ErrIsolatedNodes:
		return "Nodes are completely isolated from the graph"
	case ErrDuplicateEdge:
		return "Graph contains duplicate edges"
	case ErrOrphanedNodes:
		return "Nodes cannot reach the end"
	default:
		return "Unknown validation error"
	}
}

func setupRoutes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/validate", validateHandler)
	return mux
}

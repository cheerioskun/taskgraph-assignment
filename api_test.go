package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateHandler(t *testing.T) {
	tests := []struct {
		name           string
		request        ValidationRequest
		expectedValid  bool
		expectedErrors []string
	}{
		{
			name: "valid graph",
			request: ValidationRequest{
				Nodes: []TaskNode{
					{ID: "start", Type: StartNode},
					{ID: "middle", Type: AgentNode},
					{ID: "end", Type: EndNode},
				},
				Edges: []TaskEdge{
					{From: "start", To: "middle"},
					{From: "middle", To: "end"},
				},
			},
			expectedValid:  true,
			expectedErrors: []string{},
		},
		{
			name: "missing start node",
			request: ValidationRequest{
				Nodes: []TaskNode{
					{ID: "middle", Type: AgentNode},
					{ID: "end", Type: EndNode},
				},
				Edges: []TaskEdge{
					{From: "middle", To: "end"},
				},
			},
			expectedValid:  false,
			expectedErrors: []string{"missing_start"},
		},
		{
			name: "graph with connectivity issues",
			request: ValidationRequest{
				Nodes: []TaskNode{
					{ID: "start", Type: StartNode},
					{ID: "isolated", Type: AgentNode},
					{ID: "orphaned", Type: ToolNode},
					{ID: "end", Type: EndNode},
				},
				Edges: []TaskEdge{
					{From: "start", To: "orphaned"},
					{From: "start", To: "end"},
				},
			},
			expectedValid:  false,
			expectedErrors: []string{"isolated_nodes", "orphaned_nodes"},
		},
		{
			name: "duplicate edges",
			request: ValidationRequest{
				Nodes: []TaskNode{
					{ID: "start", Type: StartNode},
					{ID: "end", Type: EndNode},
				},
				Edges: []TaskEdge{
					{From: "start", To: "end"},
					{From: "start", To: "end"},
				},
			},
			expectedValid:  false,
			expectedErrors: []string{"duplicate_edge"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqBody, err := json.Marshal(tt.request)
			assert.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/validate", bytes.NewReader(reqBody))
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(validateHandler)
			handler.ServeHTTP(rr, req)

			assert.Equal(t, http.StatusOK, rr.Code)
			assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))

			var response ValidationResponse
			err = json.Unmarshal(rr.Body.Bytes(), &response)
			assert.NoError(t, err)

			assert.Equal(t, tt.expectedValid, response.Valid)
			assert.Len(t, response.Errors, len(tt.expectedErrors))

			actualErrorTypes := make([]string, len(response.Errors))
			for i, err := range response.Errors {
				actualErrorTypes[i] = err.Type
			}

			for _, expectedType := range tt.expectedErrors {
				assert.Contains(t, actualErrorTypes, expectedType)
			}
		})
	}
}

func TestValidateHandlerInvalidMethod(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/validate", nil)
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(validateHandler)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusMethodNotAllowed, rr.Code)
}

func TestValidateHandlerInvalidJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/validate", bytes.NewReader([]byte("invalid json")))
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(validateHandler)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

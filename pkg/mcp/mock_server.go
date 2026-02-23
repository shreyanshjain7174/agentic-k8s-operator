/*
Copyright 2026.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package mcp

import (
	"encoding/json"
	"net/http"
)

// MockServer is a mock MCP server for testing (tool-agnostic)
type MockServer struct {
	server *http.Server
}

// NewMockServer creates a new mock MCP server listening on the given address
func NewMockServer(addr string) *MockServer {
	mux := http.NewServeMux()

	// /tools endpoint returns list of available tools
	mux.HandleFunc("/tools", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		toolList := ToolListResponse{
			Tools: []string{
				"get_status",
				"propose_action",
				"execute_action",
				"validate_action",
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(toolList)
	})

	// /call_tool endpoint handles tool calls
	mux.HandleFunc("/call_tool", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req ToolRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		// Generate mock responses based on tool name
		var result map[string]interface{}
		switch req.Tool {
		case "get_status":
			result = map[string]interface{}{
				"status":    "healthy",
				"timestamp": "2026-02-23T20:00:00Z",
				"metrics": map[string]interface{}{
					"cpu_usage":    0.42,
					"memory_usage": 0.65,
					"request_rate": 1234,
				},
			}

		case "propose_action":
			result = map[string]interface{}{
				"action":      "optimize_resources",
				"description": "Reduce CPU request limits based on observed usage",
				"confidence":  "0.87",
				"impact":      "low",
			}

		case "execute_action":
			result = map[string]interface{}{
				"executed":  true,
				"action":    "optimize_resources",
				"timestamp": "2026-02-23T20:00:30Z",
				"result":    "Successfully optimized resource limits",
			}

		case "validate_action":
			result = map[string]interface{}{
				"valid":      true,
				"checks":     []string{"syntax", "permissions", "safety"},
				"violations": []string{},
			}

		default:
			w.Header().Set("Content-Type", "application/json")
			resp := ToolResponse{
				Tool:    req.Tool,
				Success: false,
				Error:   "Unknown tool: " + req.Tool,
			}
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(resp)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		resp := ToolResponse{
			Tool:    req.Tool,
			Result:  result,
			Success: true,
		}
		json.NewEncoder(w).Encode(resp)
	})

	return &MockServer{
		server: &http.Server{
			Addr:    addr,
			Handler: mux,
		},
	}
}

// Start starts the mock server (blocks until stopped)
func (ms *MockServer) Start() error {
	return ms.server.ListenAndServe()
}

// Stop stops the mock server
func (ms *MockServer) Stop() error {
	return ms.server.Close()
}

package server

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestServer_HandleRequest(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		path           string
		body           interface{}
		formatType     string
		verbose        bool
		expectedCode   int
		expectJSON     bool
		expectContains []string
	}{
		{
			name:           "GET request without body",
			method:         "GET",
			path:           "/test",
			formatType:     "go",
			verbose:        false,
			expectedCode:   http.StatusOK,
			expectJSON:     true,
			expectContains: []string{`"message": "Request processed successfully"`, `"method": "GET"`, `"path": "/test"`},
		},
		{
			name:         "POST request with JSON body - Go format",
			method:       "POST",
			path:         "/api/data",
			body:         map[string]interface{}{"name": "test", "value": 123},
			formatType:   "go",
			verbose:      true,
			expectedCode: http.StatusOK,
			expectJSON:   true,
			expectContains: []string{
				`"formatted": "type GeneratedStruct struct {`,
				"name string",
				"value float64",
			},
		},
		{
			name:         "POST request with JSON body - Rust format",
			method:       "POST",
			path:         "/api/data",
			body:         map[string]interface{}{"name": "test", "value": 123},
			formatType:   "rust",
			verbose:      true,
			expectedCode: http.StatusOK,
			expectJSON:   true,
			expectContains: []string{
				`"formatted": "#[derive(Debug, Serialize, Deserialize)]`,
				`#[serde(rename = \"name\")]`,
				`name: String`,
				`#[serde(rename = \"value\")]`,
				`value: f64`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new server instance for each test
			srv := New(8080, tt.verbose, tt.formatType)

			// Create a request
			var bodyReader *bytes.Reader
			if tt.body != nil {
				bodyBytes, err := json.Marshal(tt.body)
				if err != nil {
					t.Fatalf("Failed to marshal request body: %v", err)
				}
				bodyReader = bytes.NewReader(bodyBytes)
			} else {
				bodyReader = bytes.NewReader([]byte{})
			}

			req := httptest.NewRequest(tt.method, tt.path, bodyReader)
			if tt.body != nil {
				req.Header.Set("Content-Type", "application/json")
			}

			// Create a ResponseRecorder to record the response
			rr := httptest.NewRecorder()

			// Handle the request
			srv.handleRequest(rr, req)

			// Check status code
			if status := rr.Code; status != tt.expectedCode {
				t.Errorf("Handler returned wrong status code: got %v want %v", status, tt.expectedCode)
			}

			// Check response content
			if tt.expectJSON {
				contentType := rr.Header().Get("Content-Type")
				if !strings.Contains(contentType, "application/json") {
					t.Errorf("Handler returned wrong content type: got %v want application/json", contentType)
				}
			}

			// Check response body contains expected strings
			responseBody := rr.Body.String()
			for _, expect := range tt.expectContains {
				if !strings.Contains(responseBody, expect) {
					t.Errorf("Response body does not contain expected string: %s\nGot: %s", expect, responseBody)
				}
			}

			// For debugging
			if t.Failed() {
				t.Logf("Response body: %s", responseBody)
			}
		})
	}
}

func TestServer_Start(t *testing.T) {
	srv := New(0, false, "go") // Use port 0 to let the system assign a free port

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	// Start server in a goroutine
	errChan := make(chan error, 1)
	go func() {
		errChan <- srv.Start(ctx)
	}()

	// Wait for context cancellation or error
	select {
	case err := <-errChan:
		if err != nil && err != http.ErrServerClosed {
			t.Errorf("Server.Start() error = %v", err)
		}
	case <-ctx.Done():
		// Expected case - context timeout
	}
}

func TestFormatData(t *testing.T) {
	testData := map[string]interface{}{
		"string_field": "test",
		"number_field": 123.45,
		"bool_field":   true,
		"array_field":  []interface{}{1, 2, 3},
		"object_field": map[string]interface{}{
			"nested": "value",
		},
	}

	tests := []struct {
		name           string
		formatType     string
		expectContains []string
	}{
		{
			name:       "Go format",
			formatType: "go",
			expectContains: []string{
				"type GeneratedStruct struct",
				"string_field string",
				"number_field float64",
				"bool_field bool",
			},
		},
		{
			name:       "Rust format",
			formatType: "rust",
			expectContains: []string{
				"#[derive(Debug, Serialize, Deserialize)]",
				"struct GeneratedStruct",
				"string_field: String",
				"number_field: f64",
				"bool_field: bool",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := New(8080, false, tt.formatType)
			result, err := srv.formatData(testData)
			if err != nil {
				t.Errorf("formatData() error = %v", err)
				return
			}

			for _, expect := range tt.expectContains {
				if !strings.Contains(result, expect) {
					t.Errorf("formatData() result does not contain expected string: %s\nGot: %s", expect, result)
				}
			}
		})
	}
}

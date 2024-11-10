package server

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
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
		pretty         bool
		headers        bool
		expectedCode   int
		expectJSON     bool
		expectContains []string
		expectLogs     []string
	}{
		{
			name:           "GET request without body",
			method:         "GET",
			path:           "/test",
			formatType:     "",
			pretty:         false,
			headers:        false,
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
			pretty:       false,
			headers:      false,
			expectedCode: http.StatusOK,
			expectJSON:   true,
			expectContains: []string{
				`"message": "Request processed successfully"`,
				`"method": "POST"`,
				`"path": "/api/data"`,
			},
			expectLogs: []string{
				"JSON-Body: {",
				`"name":"test"`,
				`"value":123`,
				"type GeneratedStruct struct {",
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
			pretty:       false,
			headers:      false,
			expectedCode: http.StatusOK,
			expectJSON:   true,
			expectContains: []string{
				`"message": "Request processed successfully"`,
				`"method": "POST"`,
				`"path": "/api/data"`,
			},
			expectLogs: []string{
				"JSON-Body: {",
				`"name":"test"`,
				`"value":123`,
				"#[derive(Debug, Serialize, Deserialize)]",
				"struct GeneratedStruct",
				"#[serde(rename = \"name\")]",
				"name: String",
				"#[serde(rename = \"value\")]",
				"value: f64",
			},
		},
		{
			name:         "POST request with JSON body - Pretty Print",
			method:       "POST",
			path:         "/api/data",
			body:         map[string]interface{}{"name": "test", "value": 123},
			formatType:   "",
			pretty:       true,
			headers:      false,
			expectedCode: http.StatusOK,
			expectJSON:   true,
			expectContains: []string{
				`"message": "Request processed successfully"`,
				`"method": "POST"`,
				`"path": "/api/data"`,
			},
			expectLogs: []string{
				"JSON START",
				"JSON END",
				`"name": "test"`,
				`"value": 123`,
			},
		},
		{
			name:         "POST request with headers",
			method:       "POST",
			path:         "/api/data",
			body:         map[string]interface{}{"name": "test", "value": 123},
			formatType:   "",
			pretty:       false,
			headers:      true,
			expectedCode: http.StatusOK,
			expectJSON:   true,
			expectContains: []string{
				`"message": "Request processed successfully"`,
				`"method": "POST"`,
				`"path": "/api/data"`,
			},
			expectLogs: []string{
				"Content-Type: application/json",
				"Host: example.com",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture log output
			var logBuf bytes.Buffer
			log.SetOutput(&logBuf)
			// Restore the default logger output when the test completes
			defer log.SetOutput(os.Stderr)

			// Create a new server instance for each test
			srv := New(8080, tt.formatType, tt.pretty, tt.headers)

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

			// Check logs contain expected strings
			logOutput := logBuf.String()
			for _, expect := range tt.expectLogs {
				if !strings.Contains(logOutput, expect) {
					t.Errorf("Log output does not contain expected string: %s\nGot: %s", expect, logOutput)
				}
			}
		})
	}
}

func TestServer_Start(t *testing.T) {
	srv := New(0, "go", false, false) // Use port 0 to let the system assign a free port

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

func TestFormatJSON(t *testing.T) {
	testData := map[string]interface{}{
		"name":  "test",
		"value": 123,
	}

	tests := []struct {
		name           string
		pretty         bool
		expectContains []string
	}{
		{
			name:   "Compact JSON",
			pretty: false,
			expectContains: []string{
				"JSON-Body:",
				`"name":"test"`,
				`"value":123`,
			},
		},
		{
			name:   "Pretty JSON",
			pretty: true,
			expectContains: []string{
				"JSON START",
				"JSON END",
				`"name": "test"`,
				`"value": 123`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := New(8080, "", tt.pretty, false)
			result := srv.formatJSON(testData)

			for _, expect := range tt.expectContains {
				if !strings.Contains(result, expect) {
					t.Errorf("formatJSON() result does not contain expected string: %s\nGot: %s", expect, result)
				}
			}
		})
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
			srv := New(8080, tt.formatType, false, false)
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

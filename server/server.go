package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
)

type Server struct {
	port       int
	verbose    bool
	formatType string
}

func New(port int, verbose bool, formatType string) *Server {
	return &Server{
		port:       port,
		verbose:    verbose,
		formatType: formatType,
	}
}

func (s *Server) Start(ctx context.Context) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.handleRequest)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", s.port),
		Handler: mux,
	}

	go func() {
		<-ctx.Done()
		if err := server.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down server: %v", err)
		}
	}()

	return server.ListenAndServe()
}

func (s *Server) handleRequest(w http.ResponseWriter, r *http.Request) {
	// Log raw request if verbose mode is enabled
	if s.verbose {
		if rawRequest, err := httputil.DumpRequest(r, true); err == nil {
			log.Printf("Raw Request:\n%s", string(rawRequest))
		}
	}

	// Always log the method
	log.Printf("Received %s request to %s", r.Method, r.URL.Path)

	// Parse JSON body if present
	var bodyData interface{}
	var formatted string
	if r.Header.Get("Content-Type") == "application/json" {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		if len(body) > 0 {
			if err := json.Unmarshal(body, &bodyData); err != nil {
				http.Error(w, "Error parsing JSON", http.StatusBadRequest)
				return
			}

			// Convert and output the data in the specified format
			var err error
			formatted, err = s.formatData(bodyData)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error formatting data: %v", err), http.StatusInternalServerError)
				return
			}

			if s.verbose {
				log.Printf("Formatted data:\n%s", formatted)
			}
		}
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"message":   "Request processed successfully",
		"method":    r.Method,
		"path":      r.URL.Path,
		"formatted": formatted,
	}
	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "    ")
	encoder.Encode(response)
}

func (s *Server) formatData(data interface{}) (string, error) {
	switch s.formatType {
	case "go":
		return s.formatAsGo(data)
	case "rust":
		return s.formatAsRust(data)
	default:
		return "", fmt.Errorf("unsupported format type: %s", s.formatType)
	}
}

func (s *Server) formatAsGo(data interface{}) (string, error) {
	// Create Go struct representation
	return fmt.Sprintf("type GeneratedStruct struct {\n%s}", s.generateGoFields(data)), nil
}

func (s *Server) generateGoFields(data interface{}) string {
	switch v := data.(type) {
	case map[string]interface{}:
		var result string
		for key, val := range v {
			fieldType := s.getGoType(val)
			result += fmt.Sprintf("    %s %s `json:\"%s\"`\n", key, fieldType, key)
		}
		return result
	default:
		return "    Data interface{} `json:\"data\"`\n"
	}
}

func (s *Server) getGoType(v interface{}) string {
	switch v.(type) {
	case bool:
		return "bool"
	case float64:
		return "float64"
	case string:
		return "string"
	case []interface{}:
		return "[]interface{}"
	case map[string]interface{}:
		return "map[string]interface{}"
	case nil:
		return "interface{}"
	default:
		return "interface{}"
	}
}

func (s *Server) formatAsRust(data interface{}) (string, error) {
	// Create Rust struct representation
	return fmt.Sprintf("#[derive(Debug, Serialize, Deserialize)]\nstruct GeneratedStruct {\n%s}", s.generateRustFields(data)), nil
}

func (s *Server) generateRustFields(data interface{}) string {
	switch v := data.(type) {
	case map[string]interface{}:
		var result string
		for key, val := range v {
			fieldType := s.getRustType(val)
			result += fmt.Sprintf("    #[serde(rename = \"%s\")]\n    %s: %s,\n", key, key, fieldType)
		}
		return result
	default:
		return "    data: serde_json::Value,\n"
	}
}

func (s *Server) getRustType(v interface{}) string {
	switch v.(type) {
	case bool:
		return "bool"
	case float64:
		return "f64"
	case string:
		return "String"
	case []interface{}:
		return "Vec<serde_json::Value>"
	case map[string]interface{}:
		return "serde_json::Map<String, serde_json::Value>"
	case nil:
		return "Option<serde_json::Value>"
	default:
		return "serde_json::Value"
	}
}

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/lukehinds/reqparser/server"
)

var (
	port        = flag.Int("port", 8080, "Port to run the server on")
	formatType  = flag.String("format", "", "Output format type (go, rust) - if not provided, no struct will be generated")
	pretty      = flag.Bool("pretty", false, "Pretty print JSON with delimiters")
	headers     = flag.Bool("headers", false, "Show HTTP headers in output")
	showVersion = flag.Bool("version", false, "Show version information")
)

const version = "0.1.0"

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of reqparser:\n")
		fmt.Fprintf(os.Stderr, "\nreqparser is a HTTP request parsing and formatting tool\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		fmt.Fprintf(os.Stderr, "  -port int\n")
		fmt.Fprintf(os.Stderr, "        Port to run the server on (default 8080)\n")
		fmt.Fprintf(os.Stderr, "  -format string\n")
		fmt.Fprintf(os.Stderr, "        Output format type (go, rust) - if not provided, no struct will be generated\n")
		fmt.Fprintf(os.Stderr, "  -pretty\n")
		fmt.Fprintf(os.Stderr, "        Pretty print JSON with delimiters (if not provided, shows compact JSON-Body)\n")
		fmt.Fprintf(os.Stderr, "  -headers\n")
		fmt.Fprintf(os.Stderr, "        Show HTTP headers in output\n")
		fmt.Fprintf(os.Stderr, "  -version\n")
		fmt.Fprintf(os.Stderr, "        Show version information\n")
		fmt.Fprintf(os.Stderr, "\nBehavior:\n")
		fmt.Fprintf(os.Stderr, "  - Without -format: Shows only JSON (pretty or compact)\n")
		fmt.Fprintf(os.Stderr, "  - With -format: Shows struct and JSON (pretty or compact)\n")
		fmt.Fprintf(os.Stderr, "  - With -pretty: Shows JSON with delimiters\n")
		fmt.Fprintf(os.Stderr, "  - Without -pretty: Shows compact JSON-Body\n")
		fmt.Fprintf(os.Stderr, "  - With -headers: Shows HTTP headers\n")
		fmt.Fprintf(os.Stderr, "  - Without -headers: Headers are hidden\n")
	}

	flag.Parse()

	if *showVersion {
		fmt.Printf("reqparser version %s\n", version)
		return
	}

	// Validate format type if provided
	if *formatType != "" {
		validFormats := map[string]bool{
			"go":   true,
			"rust": true,
		}

		if !validFormats[*formatType] {
			log.Fatalf("Invalid format type: %s. Valid formats are: go, rust", *formatType)
		}
	}

	// Create server instance
	srv := server.New(*port, *formatType, *pretty, *headers)

	// Setup context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		log.Println("Shutting down server...")
		cancel()
	}()

	log.Printf("Starting server on port %d...", *port)
	if *formatType != "" {
		log.Printf("Format type: %s", *formatType)
	}
	if *pretty {
		log.Printf("Pretty JSON printing enabled")
	}
	if *headers {
		log.Printf("HTTP headers display enabled")
	}

	if err := srv.Start(ctx); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

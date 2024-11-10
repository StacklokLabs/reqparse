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
	verbose     = flag.Bool("verbose", false, "Enable verbose logging")
	formatType  = flag.String("format", "go", "Output format type (go, rust)")
	showVersion = flag.Bool("version", false, "Show version information")
)

const version = "0.1.0"

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of reqparser:\n")
		fmt.Fprintf(os.Stderr, "\nreqparser is a HTTP request parsing and formatting tool that converts JSON data into Go or Rust structs\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	if *showVersion {
		fmt.Printf("reqparser version %s\n", version)
		return
	}

	// Validate format type
	validFormats := map[string]bool{
		"go":   true,
		"rust": true,
	}

	if !validFormats[*formatType] {
		log.Fatalf("Invalid format type: %s. Valid formats are: go, rust", *formatType)
	}

	// Create server instance
	srv := server.New(*port, *verbose, *formatType)

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
	log.Printf("Format type: %s", *formatType)
	if *verbose {
		log.Printf("Verbose logging enabled")
	}

	if err := srv.Start(ctx); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

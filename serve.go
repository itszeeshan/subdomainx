package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/itszeeshan/subdomainx/internal/server"
	"github.com/itszeeshan/subdomainx/internal/tui"
)

// runServer parses server-specific flags and starts the REST API server.
func runServer(args []string) {
	fs := flag.NewFlagSet("serve", flag.ExitOnError)
	port := fs.Int("port", 8080, "Port to listen on")
	apiKey := fs.String("api-key", "", "API key for authentication (recommended)")
	outputDir := fs.String("output", "output", "Output directory for scan results")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: subdomainx serve [options]\n\nOptions:\n")
		fs.PrintDefaults()
	}

	if err := fs.Parse(args); err != nil {
		log.Fatalf("Failed to parse serve flags: %v", err)
	}

	// Register the scan pipeline function so the server can invoke it.
	tui.RegisterScanFunc(tuiScanFunc)

	srv := server.New(*port, *apiKey, *outputDir)
	if err := srv.Start(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

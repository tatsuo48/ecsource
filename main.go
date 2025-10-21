package main

import (
	"fmt"
	"log"
	"os"

	"github.com/tatsuo48/ecsource/internal/calculator"
	"github.com/tatsuo48/ecsource/internal/parser"
	"github.com/tatsuo48/ecsource/internal/renderer"
)

// These variables are set in build step
var (
	Version  = "unset"
	Revision = "unset"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("please specify file path")
	}

	if os.Args[1] == "-h" || os.Args[1] == "--help" {
		fmt.Printf("Usage: %s [file path]", os.Args[0])
		os.Exit(0)
	}

	if os.Args[1] == "--version" || os.Args[1] == "-v" {
		fmt.Printf("version: %s\n", Version)
		fmt.Printf("revision: %s\n", Revision)
		os.Exit(0)
	}

	if err := run(os.Args[1]); err != nil {
		log.Fatal(err)
	}
}

// run executes the main logic of loading, calculating, and displaying task definition resources
func run(filepath string) error {
	// Load and parse task definition
	taskDef, err := parser.LoadTaskDefinitionFromFile(filepath)
	if err != nil {
		return fmt.Errorf("failed to load task definition: %w", err)
	}

	// Calculate resource summary
	summary, err := calculator.CalculateResources(taskDef)
	if err != nil {
		return fmt.Errorf("failed to calculate resources: %w", err)
	}

	// Render the table
	renderer.RenderTable(os.Stdout, taskDef, summary)

	return nil
}

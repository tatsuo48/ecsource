package main

import (
	"flag"
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

// CLI flags
var (
	filepath     string
	outputFormat string
	showVersion  bool
	noColor      bool
)

func init() {
	flag.StringVar(&filepath, "f", "", "Path to ECS task definition JSON file")
	flag.StringVar(&filepath, "file", "", "Path to ECS task definition JSON file (alias)")
	flag.StringVar(&outputFormat, "o", "table", "Output format: table, json, csv")
	flag.StringVar(&outputFormat, "output", "table", "Output format: table, json, csv (alias)")
	flag.BoolVar(&showVersion, "v", false, "Show version information")
	flag.BoolVar(&showVersion, "version", false, "Show version information (alias)")
	flag.BoolVar(&noColor, "no-color", false, "Disable color output")

	flag.Usage = printUsage
}

func printUsage() {
	fmt.Fprintf(os.Stderr, `ecsource - ECS Task Definition resource analyzer

Usage:
  %s [OPTIONS] [FILE]
  %s -f FILE [OPTIONS]

Options:
  -f, --file string       Path to ECS task definition JSON file
  -o, --output string     Output format: table, json, csv (default: table)
  -v, --version           Show version information
  --no-color              Disable color output
  -h, --help              Show this help message

Examples:
  %s task.json
  %s -f task.json -o json
  %s --file task.json --output csv --no-color

`, os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0])
}

func printVersion() {
	fmt.Printf("ecsource version %s\n", Version)
	fmt.Printf("revision: %s\n", Revision)
}

func main() {
	flag.Parse()

	// Handle version flag
	if showVersion {
		printVersion()
		os.Exit(0)
	}

	// Determine file path from flag or positional argument
	if filepath == "" {
		// Try to get from positional argument
		if flag.NArg() > 0 {
			filepath = flag.Arg(0)
		} else {
			fmt.Fprintf(os.Stderr, "Error: No file specified\n\n")
			flag.Usage()
			os.Exit(1)
		}
	}

	// Validate output format
	validFormats := map[string]bool{
		"table": true,
		"json":  true,
		"csv":   true,
	}
	if !validFormats[outputFormat] {
		fmt.Fprintf(os.Stderr, "Error: Invalid output format '%s'. Valid formats: table, json, csv\n", outputFormat)
		os.Exit(1)
	}

	// Run the main logic
	if err := run(filepath, outputFormat, noColor); err != nil {
		log.Fatal(err)
	}
}

// run executes the main logic of loading, calculating, and displaying task definition resources
func run(filepath, format string, noColor bool) error {
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

	// Render output based on format
	switch format {
	case "table":
		renderer.RenderTable(os.Stdout, taskDef, summary, !noColor)
	case "json":
		renderer.RenderJSON(os.Stdout, taskDef, summary)
	case "csv":
		renderer.RenderCSV(os.Stdout, taskDef, summary)
	default:
		return fmt.Errorf("unsupported output format: %s", format)
	}

	return nil
}

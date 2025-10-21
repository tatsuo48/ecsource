# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

ecsource is a CLI tool written in Go that parses AWS ECS Task Definition JSON files and displays resource allocations (CPU, memory, memory reservations) in multiple output formats. It calculates and highlights leftover resources, displaying negative values in red (table format) to indicate over-allocation.

## Development Commands

### Build
```bash
make build
# or
go build -o bin/ecsource main.go
```

### Install locally
```bash
make install
```

### Clean
```bash
make clean
```

### Run the tool
```bash
# Table output (default)
./bin/ecsource <path-to-task-definition.json>
ecsource test_task.json

# JSON output
ecsource -f test_task.json -o json

# CSV output
ecsource -f test_task.json -o csv

# Disable color output
ecsource test_task.json --no-color

# Show version
ecsource --version

# Show help
ecsource --help
```

### Testing
```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with verbose output
go test -v ./...

# Run tests with race detection
go test -race ./...
```

## Architecture

The project follows a modular package-based architecture:

### Package Structure
```
ecsource/
├── main.go                      # CLI entry point and argument parsing
├── internal/
│   ├── models/                  # Data structures
│   │   └── task.go             # TaskDefinition, ContainerDefinition, ResourceSummary
│   ├── parser/                  # JSON parsing logic
│   │   ├── parser.go           # LoadTaskDefinitionFromFile, ParseTaskDefinition
│   │   └── parser_test.go      # Parser tests (92.3% coverage)
│   ├── calculator/              # Resource calculation logic
│   │   ├── calculator.go       # CalculateResources
│   │   └── calculator_test.go  # Calculator tests (100% coverage)
│   └── renderer/                # Output rendering logic
│       ├── table.go            # RenderTable, RenderJSON, RenderCSV
│       └── table_test.go       # Renderer tests (100% coverage)
```

### Input Format
The tool accepts ECS Task Definition JSON in two formats:
1. Wrapped format: `{"taskDefinition": {...}}`
2. Direct format: `{...}` (raw task definition object)

The JSON is parsed by `internal/parser` package into the `TaskDefinition` struct.

### Core Data Structures (internal/models)
- `TaskDefinition`: Top-level structure representing the ECS task
- `ContainerDefinition`: Individual container configuration including CPU, Memory, MemoryReservation
- `ResourceSummary`: Calculated resource totals and leftovers
- Supporting structs: `Environment`, `LogConfiguration`, `Secret`, `PortMapping`, `FirelensConfiguration`, `DockerLabels`

### Processing Flow
1. **CLI Parsing** (`main.go`): Parse command-line arguments using flag package
2. **File Loading** (`internal/parser`): Read and parse JSON file into TaskDefinition
3. **Calculation** (`internal/calculator`): Calculate resource totals and leftovers
4. **Rendering** (`internal/renderer`): Output results in specified format (table/JSON/CSV)

### CLI Options
- `-f, --file`: Path to ECS task definition JSON file
- `-o, --output`: Output format (table, json, csv) - default: table
- `-v, --version`: Show version information
- `--no-color`: Disable color output for table format
- `-h, --help`: Show help message

### Version Information
The `Version` and `Revision` variables are set at build time via ldflags (see `.goreleaser.yml`). These are displayed with the `-v` or `--version` flag.

## Release Process

This project uses automated release management:
- **tagpr**: Automatically creates release PRs when changes are merged to main
- **goreleaser**: Builds binaries and publishes releases when tags are created
- Homebrew formula is automatically updated in `tatsuo48/homebrew-tap`

Version bumps and changelog updates are handled automatically by tagpr.

## Dependencies

- `github.com/olekukonko/tablewriter`: For rendering ASCII tables
- Go 1.19+
- Standard library: encoding/json, encoding/csv, flag, io, fmt, log, os, strconv

## Key Implementation Details

### Output Formats

#### Table Format (default)
- Displays task-level settings, individual container resources, sum of all containers, and leftover resources
- Negative leftover values (over-allocation) are highlighted in red using tablewriter's color feature
- Color can be disabled with `--no-color` flag

#### JSON Format
- Structured JSON output with taskDefinition, containers, and summary sections
- Useful for programmatic processing and integration with other tools

#### CSV Format
- Simple CSV format with header row
- Easy to import into spreadsheets or data analysis tools

### Error Handling
- All parsing errors are properly handled and returned with context
- Invalid CPU/Memory values are caught and reported with clear error messages
- File not found and invalid JSON errors are handled gracefully

### Testing
- Comprehensive test suite with 25+ test cases
- Coverage: main (50%+), internal packages (92-100%)
- Tests cover normal operation, error cases, and edge cases
- CI/CD integration with GitHub Actions for automated testing

### Code Quality
- All code formatted with `gofmt`
- No warnings from `go vet`
- Clear separation of concerns following SOLID principles
- Well-documented public APIs with GoDoc comments

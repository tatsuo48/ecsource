# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

ecsource is a CLI tool written in Go that parses AWS ECS Task Definition JSON files and displays resource allocations (CPU, memory, memory reservations) in an easy-to-read table format. It calculates and highlights leftover resources, displaying negative values in red to indicate over-allocation.

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
./bin/ecsource <path-to-task-definition.json>
# or
ecsource test_task.json
```

### Testing
Currently, there are no automated tests in this project. When adding tests, follow Go conventions:
- Create `*_test.go` files
- Run tests with `go test ./...`

## Architecture

This is a single-file CLI application (`main.go`) with a straightforward structure:

### Input Format
The tool accepts ECS Task Definition JSON in two formats:
1. Wrapped format: `{"taskDefinition": {...}}`
2. Direct format: `{...}` (raw task definition object)

The JSON is unmarshalled into the `TaskDefinition` struct, which contains:
- Task-level settings: CPU, Memory, Family, NetworkMode, etc.
- Array of `ContainerDefinition` objects with per-container resources

### Core Data Structures
- `TaskDefinition`: Top-level structure representing the ECS task
- `ContainerDefinition`: Individual container configuration including CPU, Memory, MemoryReservation
- Supporting structs: `Environment`, `LogConfiguration`, `Secret`, `PortMapping`, `FirelensConfiguration`, `DockerLabels`

### Processing Logic
1. Parse command-line arguments (file path, version, help)
2. Read and unmarshal JSON file into `TaskDefinition`
3. Create table with tablewriter library
4. Calculate sums of all container resources
5. Calculate leftover resources (task allocation - sum of containers)
6. Apply red color to negative leftover values
7. Render table to stdout

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

## Key Implementation Details

- The tool displays task-level settings, individual container resources, sum of all containers, and leftover resources
- Negative leftover values (over-allocation) are highlighted in red using tablewriter's color feature (main.go:159-176)
- The tool handles both wrapped (`{"taskDefinition": {...}}`) and unwrapped JSON formats for flexibility
- Error handling uses `log.Fatal()` for all errors

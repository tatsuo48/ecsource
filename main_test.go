package main

import (
	"os"
	"testing"
)

func TestRunWithValidFile(t *testing.T) {
	// Check if test_task.json exists
	if _, err := os.Stat("test_task.json"); os.IsNotExist(err) {
		t.Skip("test_task.json not found, skipping test")
	}

	err := run("test_task.json", "table", false)
	if err != nil {
		t.Fatalf("Expected run() to succeed with test_task.json, got error: %v", err)
	}
}

func TestRunWithNonexistentFile(t *testing.T) {
	err := run("nonexistent_file.json", "table", false)
	if err == nil {
		t.Error("Expected error for nonexistent file, got nil")
	}
}

func TestRunWithInvalidJSON(t *testing.T) {
	// Create a temporary file with invalid JSON
	tmpFile, err := os.CreateTemp("", "invalid_*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.Write([]byte("{invalid json}")); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	err = run(tmpFile.Name(), "table", false)
	if err == nil {
		t.Error("Expected error for invalid JSON, got nil")
	}
}

func TestRunWithInvalidCPU(t *testing.T) {
	// Create a temporary file with invalid CPU value
	tmpFile, err := os.CreateTemp("", "invalid_cpu_*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	invalidJSON := `{
		"family": "test",
		"cpu": "invalid",
		"memory": "1024",
		"containerDefinitions": [
			{
				"name": "test",
				"image": "test:latest",
				"cpu": 256,
				"memory": 512,
				"memoryReservation": 256,
				"essential": true,
				"portMappings": [],
				"environment": [],
				"logConfiguration": {"logDriver": "awslogs"},
				"secrets": [],
				"command": [],
				"volumesFrom": []
			}
		],
		"placementConstraints": []
	}`

	if _, err := tmpFile.Write([]byte(invalidJSON)); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	err = run(tmpFile.Name(), "table", false)
	if err == nil {
		t.Error("Expected error for invalid CPU value, got nil")
	}
}

func TestRunWithJSONOutput(t *testing.T) {
	if _, err := os.Stat("test_task.json"); os.IsNotExist(err) {
		t.Skip("test_task.json not found, skipping test")
	}

	err := run("test_task.json", "json", false)
	if err != nil {
		t.Fatalf("Expected run() to succeed with JSON output, got error: %v", err)
	}
}

func TestRunWithCSVOutput(t *testing.T) {
	if _, err := os.Stat("test_task.json"); os.IsNotExist(err) {
		t.Skip("test_task.json not found, skipping test")
	}

	err := run("test_task.json", "csv", false)
	if err != nil {
		t.Fatalf("Expected run() to succeed with CSV output, got error: %v", err)
	}
}

func TestRunWithInvalidFormat(t *testing.T) {
	if _, err := os.Stat("test_task.json"); os.IsNotExist(err) {
		t.Skip("test_task.json not found, skipping test")
	}

	err := run("test_task.json", "invalid", false)
	if err == nil {
		t.Error("Expected error for invalid output format, got nil")
	}
}

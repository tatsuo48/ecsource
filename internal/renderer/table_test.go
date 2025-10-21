package renderer

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/tatsuo48/ecsource/internal/models"
)

func TestRenderTable(t *testing.T) {
	taskDef := &models.TaskDefinition{
		CPU:    "1024",
		Memory: "2048",
		ContainerDefinitions: []models.ContainerDefinition{
			{
				Name:              "container1",
				CPU:               256,
				Memory:            512,
				MemoryReservation: 256,
			},
			{
				Name:              "container2",
				CPU:               256,
				Memory:            512,
				MemoryReservation: 256,
			},
		},
	}

	summary := &models.ResourceSummary{
		TaskCPU:                   1024,
		TaskMemory:                2048,
		TotalCPU:                  512,
		TotalMemory:               1024,
		TotalMemoryReservation:    512,
		LeftoverCPU:               512,
		LeftoverMemory:            1024,
		LeftoverMemoryReservation: 1536,
	}

	var buf bytes.Buffer
	RenderTable(&buf, taskDef, summary, true)

	output := buf.String()

	// Verify key elements are in the output
	expectedElements := []string{
		"NAME",
		"CPU",
		"MEMORY",
		"MEMORYRESERVATION",
		"Task Setting",
		"container1",
		"container2",
		"sum of all container",
		"LEFTOVER",
		"1024",
		"2048",
		"512",
	}

	for _, elem := range expectedElements {
		if !strings.Contains(output, elem) {
			t.Errorf("Expected output to contain '%s', but it doesn't.\nOutput:\n%s", elem, output)
		}
	}
}

func TestRenderTableNegativeLeftover(t *testing.T) {
	taskDef := &models.TaskDefinition{
		CPU:    "512",
		Memory: "1024",
		ContainerDefinitions: []models.ContainerDefinition{
			{
				Name:              "overallocated",
				CPU:               600,
				Memory:            1500,
				MemoryReservation: 1200,
			},
		},
	}

	summary := &models.ResourceSummary{
		TaskCPU:                   512,
		TaskMemory:                1024,
		TotalCPU:                  600,
		TotalMemory:               1500,
		TotalMemoryReservation:    1200,
		LeftoverCPU:               -88,
		LeftoverMemory:            -476,
		LeftoverMemoryReservation: -176,
	}

	var buf bytes.Buffer
	RenderTable(&buf, taskDef, summary, true)

	output := buf.String()

	// Verify negative values are present
	// Note: We can't easily test for red color in the output without more sophisticated parsing,
	// but we can verify the negative values are displayed
	if !strings.Contains(output, "-88") {
		t.Errorf("Expected output to contain '-88', but it doesn't.\nOutput:\n%s", output)
	}
	if !strings.Contains(output, "-476") {
		t.Errorf("Expected output to contain '-476', but it doesn't.\nOutput:\n%s", output)
	}
}

func TestRenderTableEmptyContainers(t *testing.T) {
	taskDef := &models.TaskDefinition{
		CPU:                  "1024",
		Memory:               "2048",
		ContainerDefinitions: []models.ContainerDefinition{},
	}

	summary := &models.ResourceSummary{
		TaskCPU:                   1024,
		TaskMemory:                2048,
		TotalCPU:                  0,
		TotalMemory:               0,
		TotalMemoryReservation:    0,
		LeftoverCPU:               1024,
		LeftoverMemory:            2048,
		LeftoverMemoryReservation: 2048,
	}

	var buf bytes.Buffer
	RenderTable(&buf, taskDef, summary, true)

	output := buf.String()

	// Verify task settings are shown even with no containers
	if !strings.Contains(output, "Task Setting") {
		t.Errorf("Expected output to contain 'Task Setting', but it doesn't.\nOutput:\n%s", output)
	}
}

func TestRenderTableNoColor(t *testing.T) {
	taskDef := &models.TaskDefinition{
		CPU:    "512",
		Memory: "1024",
		ContainerDefinitions: []models.ContainerDefinition{
			{
				Name:              "test",
				CPU:               600,
				Memory:            1500,
				MemoryReservation: 1024,
			},
		},
	}

	summary := &models.ResourceSummary{
		TaskCPU:                   512,
		TaskMemory:                1024,
		TotalCPU:                  600,
		TotalMemory:               1500,
		TotalMemoryReservation:    1024,
		LeftoverCPU:               -88,
		LeftoverMemory:            -476,
		LeftoverMemoryReservation: 0,
	}

	var buf bytes.Buffer
	RenderTable(&buf, taskDef, summary, false)

	output := buf.String()

	// Verify output is generated
	if !strings.Contains(output, "test") {
		t.Errorf("Expected output to contain container name, but it doesn't.\nOutput:\n%s", output)
	}
}

func TestRenderJSON(t *testing.T) {
	taskDef := &models.TaskDefinition{
		Family: "test-family",
		CPU:    "1024",
		Memory: "2048",
		ContainerDefinitions: []models.ContainerDefinition{
			{
				Name:              "container1",
				CPU:               256,
				Memory:            512,
				MemoryReservation: 256,
			},
		},
	}

	summary := &models.ResourceSummary{
		TaskCPU:                   1024,
		TaskMemory:                2048,
		TotalCPU:                  256,
		TotalMemory:               512,
		TotalMemoryReservation:    256,
		LeftoverCPU:               768,
		LeftoverMemory:            1536,
		LeftoverMemoryReservation: 1792,
	}

	var buf bytes.Buffer
	err := RenderJSON(&buf, taskDef, summary)
	if err != nil {
		t.Fatalf("RenderJSON failed: %v", err)
	}

	// Parse the JSON output
	var output OutputData
	if err := json.Unmarshal(buf.Bytes(), &output); err != nil {
		t.Fatalf("Failed to parse JSON output: %v\nOutput:\n%s", err, buf.String())
	}

	// Verify task definition
	if output.TaskDefinition.Family != "test-family" {
		t.Errorf("Expected family 'test-family', got '%s'", output.TaskDefinition.Family)
	}
	if output.TaskDefinition.CPU != 1024 {
		t.Errorf("Expected CPU 1024, got %d", output.TaskDefinition.CPU)
	}

	// Verify containers
	if len(output.Containers) != 1 {
		t.Errorf("Expected 1 container, got %d", len(output.Containers))
	}
	if output.Containers[0].Name != "container1" {
		t.Errorf("Expected container name 'container1', got '%s'", output.Containers[0].Name)
	}

	// Verify summary
	if output.Summary.LeftoverCPU != 768 {
		t.Errorf("Expected leftover CPU 768, got %d", output.Summary.LeftoverCPU)
	}
}

func TestRenderCSV(t *testing.T) {
	taskDef := &models.TaskDefinition{
		Family: "test-family",
		CPU:    "1024",
		Memory: "2048",
		ContainerDefinitions: []models.ContainerDefinition{
			{
				Name:              "container1",
				CPU:               256,
				Memory:            512,
				MemoryReservation: 256,
			},
		},
	}

	summary := &models.ResourceSummary{
		TaskCPU:                   1024,
		TaskMemory:                2048,
		TotalCPU:                  256,
		TotalMemory:               512,
		TotalMemoryReservation:    256,
		LeftoverCPU:               768,
		LeftoverMemory:            1536,
		LeftoverMemoryReservation: 1792,
	}

	var buf bytes.Buffer
	err := RenderCSV(&buf, taskDef, summary)
	if err != nil {
		t.Fatalf("RenderCSV failed: %v", err)
	}

	output := buf.String()

	// Verify CSV structure
	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) != 5 {
		t.Errorf("Expected 5 CSV lines (header + task + container + sum + leftover), got %d", len(lines))
	}

	// Verify header
	if !strings.Contains(lines[0], "Name") {
		t.Errorf("Expected CSV header to contain 'Name', got: %s", lines[0])
	}

	// Verify container data
	if !strings.Contains(output, "container1") {
		t.Errorf("Expected CSV to contain 'container1', but it doesn't.\nOutput:\n%s", output)
	}

	// Verify leftover row
	if !strings.Contains(output, "leftover") {
		t.Errorf("Expected CSV to contain 'leftover', but it doesn't.\nOutput:\n%s", output)
	}
}

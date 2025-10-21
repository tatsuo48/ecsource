package renderer

import (
	"bytes"
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
	RenderTable(&buf, taskDef, summary)

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
	RenderTable(&buf, taskDef, summary)

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
	RenderTable(&buf, taskDef, summary)

	output := buf.String()

	// Verify task settings are shown even with no containers
	if !strings.Contains(output, "Task Setting") {
		t.Errorf("Expected output to contain 'Task Setting', but it doesn't.\nOutput:\n%s", output)
	}
}

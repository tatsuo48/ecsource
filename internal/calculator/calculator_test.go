package calculator

import (
	"testing"

	"github.com/tatsuo48/ecsource/internal/models"
)

func TestCalculateResources(t *testing.T) {
	tests := []struct {
		name                           string
		taskDef                        *models.TaskDefinition
		expectedTaskCPU                int64
		expectedTaskMemory             int64
		expectedTotalCPU               int64
		expectedTotalMemory            int64
		expectedTotalMemReservation    int64
		expectedLeftoverCPU            int64
		expectedLeftoverMemory         int64
		expectedLeftoverMemReservation int64
		expectError                    bool
	}{
		{
			name: "Positive leftovers",
			taskDef: &models.TaskDefinition{
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
			},
			expectedTaskCPU:                1024,
			expectedTaskMemory:             2048,
			expectedTotalCPU:               512,
			expectedTotalMemory:            1024,
			expectedTotalMemReservation:    512,
			expectedLeftoverCPU:            512,
			expectedLeftoverMemory:         1024,
			expectedLeftoverMemReservation: 1536,
			expectError:                    false,
		},
		{
			name: "Negative leftovers (over-allocation)",
			taskDef: &models.TaskDefinition{
				CPU:    "512",
				Memory: "1024",
				ContainerDefinitions: []models.ContainerDefinition{
					{
						Name:              "container1",
						CPU:               600,
						Memory:            1500,
						MemoryReservation: 1024,
					},
				},
			},
			expectedTaskCPU:                512,
			expectedTaskMemory:             1024,
			expectedTotalCPU:               600,
			expectedTotalMemory:            1500,
			expectedTotalMemReservation:    1024,
			expectedLeftoverCPU:            -88,
			expectedLeftoverMemory:         -476,
			expectedLeftoverMemReservation: 0,
			expectError:                    false,
		},
		{
			name: "Zero leftovers (exact match)",
			taskDef: &models.TaskDefinition{
				CPU:    "1024",
				Memory: "2048",
				ContainerDefinitions: []models.ContainerDefinition{
					{
						Name:              "container1",
						CPU:               1024,
						Memory:            2048,
						MemoryReservation: 2048,
					},
				},
			},
			expectedTaskCPU:                1024,
			expectedTaskMemory:             2048,
			expectedTotalCPU:               1024,
			expectedTotalMemory:            2048,
			expectedTotalMemReservation:    2048,
			expectedLeftoverCPU:            0,
			expectedLeftoverMemory:         0,
			expectedLeftoverMemReservation: 0,
			expectError:                    false,
		},
		{
			name: "Invalid CPU format",
			taskDef: &models.TaskDefinition{
				CPU:                  "invalid",
				Memory:               "1024",
				ContainerDefinitions: []models.ContainerDefinition{},
			},
			expectError: true,
		},
		{
			name: "Invalid Memory format",
			taskDef: &models.TaskDefinition{
				CPU:                  "512",
				Memory:               "invalid",
				ContainerDefinitions: []models.ContainerDefinition{},
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			summary, err := CalculateResources(tt.taskDef)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if summary.TaskCPU != tt.expectedTaskCPU {
				t.Errorf("Expected TaskCPU %d, got %d", tt.expectedTaskCPU, summary.TaskCPU)
			}
			if summary.TaskMemory != tt.expectedTaskMemory {
				t.Errorf("Expected TaskMemory %d, got %d", tt.expectedTaskMemory, summary.TaskMemory)
			}
			if summary.TotalCPU != tt.expectedTotalCPU {
				t.Errorf("Expected TotalCPU %d, got %d", tt.expectedTotalCPU, summary.TotalCPU)
			}
			if summary.TotalMemory != tt.expectedTotalMemory {
				t.Errorf("Expected TotalMemory %d, got %d", tt.expectedTotalMemory, summary.TotalMemory)
			}
			if summary.TotalMemoryReservation != tt.expectedTotalMemReservation {
				t.Errorf("Expected TotalMemoryReservation %d, got %d", tt.expectedTotalMemReservation, summary.TotalMemoryReservation)
			}
			if summary.LeftoverCPU != tt.expectedLeftoverCPU {
				t.Errorf("Expected LeftoverCPU %d, got %d", tt.expectedLeftoverCPU, summary.LeftoverCPU)
			}
			if summary.LeftoverMemory != tt.expectedLeftoverMemory {
				t.Errorf("Expected LeftoverMemory %d, got %d", tt.expectedLeftoverMemory, summary.LeftoverMemory)
			}
			if summary.LeftoverMemoryReservation != tt.expectedLeftoverMemReservation {
				t.Errorf("Expected LeftoverMemoryReservation %d, got %d", tt.expectedLeftoverMemReservation, summary.LeftoverMemoryReservation)
			}
		})
	}
}

func TestCalculateResourcesEmptyContainers(t *testing.T) {
	taskDef := &models.TaskDefinition{
		CPU:                  "1024",
		Memory:               "2048",
		ContainerDefinitions: []models.ContainerDefinition{},
	}

	summary, err := CalculateResources(taskDef)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if summary.TotalCPU != 0 {
		t.Errorf("Expected TotalCPU 0 for empty containers, got %d", summary.TotalCPU)
	}
	if summary.LeftoverCPU != 1024 {
		t.Errorf("Expected LeftoverCPU 1024, got %d", summary.LeftoverCPU)
	}
}

package calculator

import (
	"fmt"
	"strconv"

	"github.com/tatsuo48/ecsource/internal/models"
)

// CalculateResources computes resource totals and leftovers for a task definition.
// It returns a ResourceSummary with all calculated values.
func CalculateResources(taskDef *models.TaskDefinition) (*models.ResourceSummary, error) {
	// Parse task-level resources
	taskCPU, err := strconv.ParseInt(taskDef.CPU, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse task CPU '%s': %w", taskDef.CPU, err)
	}

	taskMemory, err := strconv.ParseInt(taskDef.Memory, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse task Memory '%s': %w", taskDef.Memory, err)
	}

	// Calculate totals from all containers
	var totalCPU, totalMemory, totalMemoryReservation int64
	for _, container := range taskDef.ContainerDefinitions {
		totalCPU += container.CPU
		totalMemory += container.Memory
		totalMemoryReservation += container.MemoryReservation
	}

	// Calculate leftovers
	leftoverCPU := taskCPU - totalCPU
	leftoverMemory := taskMemory - totalMemory
	leftoverMemoryReservation := taskMemory - totalMemoryReservation

	return &models.ResourceSummary{
		TaskCPU:                   taskCPU,
		TaskMemory:                taskMemory,
		TotalCPU:                  totalCPU,
		TotalMemory:               totalMemory,
		TotalMemoryReservation:    totalMemoryReservation,
		LeftoverCPU:               leftoverCPU,
		LeftoverMemory:            leftoverMemory,
		LeftoverMemoryReservation: leftoverMemoryReservation,
	}, nil
}

package renderer

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"

	"github.com/olekukonko/tablewriter"
	"github.com/tatsuo48/ecsource/internal/models"
)

// RenderTable renders a formatted table of task definition resources to the given writer.
// Negative leftover values are displayed in red to indicate over-allocation (if color is enabled).
func RenderTable(w io.Writer, taskDef *models.TaskDefinition, summary *models.ResourceSummary, enableColor bool) {
	table := tablewriter.NewWriter(w)
	table.SetHeader([]string{"Name", "CPU", "Memory", "MemoryReservation"})

	// Add task-level settings row
	table.Append([]string{
		"Task Setting",
		fmt.Sprintf("%d", summary.TaskCPU),
		fmt.Sprintf("%d", summary.TaskMemory),
		fmt.Sprintf("%d", summary.TaskMemory),
	})

	// Add container rows
	for _, container := range taskDef.ContainerDefinitions {
		table.Append([]string{
			container.Name,
			fmt.Sprintf("%d", container.CPU),
			fmt.Sprintf("%d", container.Memory),
			fmt.Sprintf("%d", container.MemoryReservation),
		})
	}

	// Add sum row
	table.Append([]string{
		"sum of all container",
		fmt.Sprintf("%d", summary.TotalCPU),
		fmt.Sprintf("%d", summary.TotalMemory),
		fmt.Sprintf("%d", summary.TotalMemoryReservation),
	})

	// Set footer with leftover values
	footerStrings := []string{
		"leftover",
		fmt.Sprintf("%d", summary.LeftoverCPU),
		fmt.Sprintf("%d", summary.LeftoverMemory),
		fmt.Sprintf("%d", summary.LeftoverMemoryReservation),
	}
	table.SetFooter(footerStrings)

	// Apply red color to negative leftover values (only if color is enabled)
	if enableColor {
		footerColors := make([]tablewriter.Colors, 4)
		footerColors[0] = tablewriter.Colors{} // "leftover" label - no color

		if summary.LeftoverCPU < 0 {
			footerColors[1] = tablewriter.Colors{tablewriter.FgRedColor}
		} else {
			footerColors[1] = tablewriter.Colors{}
		}

		if summary.LeftoverMemory < 0 {
			footerColors[2] = tablewriter.Colors{tablewriter.FgRedColor}
		} else {
			footerColors[2] = tablewriter.Colors{}
		}

		if summary.LeftoverMemoryReservation < 0 {
			footerColors[3] = tablewriter.Colors{tablewriter.FgRedColor}
		} else {
			footerColors[3] = tablewriter.Colors{}
		}

		table.SetFooterColor(footerColors...)
	}

	table.Render()
}

// OutputData represents the complete output data structure for JSON export
type OutputData struct {
	TaskDefinition TaskInfo        `json:"taskDefinition"`
	Containers     []ContainerInfo `json:"containers"`
	Summary        SummaryInfo     `json:"summary"`
}

// TaskInfo represents task-level information
type TaskInfo struct {
	Family string `json:"family"`
	CPU    int64  `json:"cpu"`
	Memory int64  `json:"memory"`
}

// ContainerInfo represents container-level information
type ContainerInfo struct {
	Name              string `json:"name"`
	CPU               int64  `json:"cpu"`
	Memory            int64  `json:"memory"`
	MemoryReservation int64  `json:"memoryReservation"`
}

// SummaryInfo represents summary calculations
type SummaryInfo struct {
	TotalCPU                  int64 `json:"totalCpu"`
	TotalMemory               int64 `json:"totalMemory"`
	TotalMemoryReservation    int64 `json:"totalMemoryReservation"`
	LeftoverCPU               int64 `json:"leftoverCpu"`
	LeftoverMemory            int64 `json:"leftoverMemory"`
	LeftoverMemoryReservation int64 `json:"leftoverMemoryReservation"`
}

// RenderJSON renders the task definition and summary as JSON
func RenderJSON(w io.Writer, taskDef *models.TaskDefinition, summary *models.ResourceSummary) error {
	// Build container info
	containers := make([]ContainerInfo, len(taskDef.ContainerDefinitions))
	for i, c := range taskDef.ContainerDefinitions {
		containers[i] = ContainerInfo{
			Name:              c.Name,
			CPU:               c.CPU,
			Memory:            c.Memory,
			MemoryReservation: c.MemoryReservation,
		}
	}

	// Build output data
	output := OutputData{
		TaskDefinition: TaskInfo{
			Family: taskDef.Family,
			CPU:    summary.TaskCPU,
			Memory: summary.TaskMemory,
		},
		Containers: containers,
		Summary: SummaryInfo{
			TotalCPU:                  summary.TotalCPU,
			TotalMemory:               summary.TotalMemory,
			TotalMemoryReservation:    summary.TotalMemoryReservation,
			LeftoverCPU:               summary.LeftoverCPU,
			LeftoverMemory:            summary.LeftoverMemory,
			LeftoverMemoryReservation: summary.LeftoverMemoryReservation,
		},
	}

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(output)
}

// RenderCSV renders the task definition and summary as CSV
func RenderCSV(w io.Writer, taskDef *models.TaskDefinition, summary *models.ResourceSummary) error {
	writer := csv.NewWriter(w)
	defer writer.Flush()

	// Write header
	if err := writer.Write([]string{"Name", "CPU", "Memory", "MemoryReservation"}); err != nil {
		return fmt.Errorf("failed to write CSV header: %w", err)
	}

	// Write task setting
	if err := writer.Write([]string{
		"Task Setting",
		fmt.Sprintf("%d", summary.TaskCPU),
		fmt.Sprintf("%d", summary.TaskMemory),
		fmt.Sprintf("%d", summary.TaskMemory),
	}); err != nil {
		return fmt.Errorf("failed to write task setting: %w", err)
	}

	// Write container rows
	for _, container := range taskDef.ContainerDefinitions {
		if err := writer.Write([]string{
			container.Name,
			fmt.Sprintf("%d", container.CPU),
			fmt.Sprintf("%d", container.Memory),
			fmt.Sprintf("%d", container.MemoryReservation),
		}); err != nil {
			return fmt.Errorf("failed to write container row: %w", err)
		}
	}

	// Write sum row
	if err := writer.Write([]string{
		"sum of all container",
		fmt.Sprintf("%d", summary.TotalCPU),
		fmt.Sprintf("%d", summary.TotalMemory),
		fmt.Sprintf("%d", summary.TotalMemoryReservation),
	}); err != nil {
		return fmt.Errorf("failed to write sum row: %w", err)
	}

	// Write leftover row
	if err := writer.Write([]string{
		"leftover",
		fmt.Sprintf("%d", summary.LeftoverCPU),
		fmt.Sprintf("%d", summary.LeftoverMemory),
		fmt.Sprintf("%d", summary.LeftoverMemoryReservation),
	}); err != nil {
		return fmt.Errorf("failed to write leftover row: %w", err)
	}

	return nil
}

package renderer

import (
	"fmt"
	"io"

	"github.com/olekukonko/tablewriter"
	"github.com/tatsuo48/ecsource/internal/models"
)

// RenderTable renders a formatted table of task definition resources to the given writer.
// Negative leftover values are displayed in red to indicate over-allocation.
func RenderTable(w io.Writer, taskDef *models.TaskDefinition, summary *models.ResourceSummary) {
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

	// Apply red color to negative leftover values
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
	table.Render()
}

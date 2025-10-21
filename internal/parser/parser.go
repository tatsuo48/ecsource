package parser

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/tatsuo48/ecsource/internal/models"
)

// LoadTaskDefinitionFromFile reads and parses an ECS task definition JSON file.
// It supports both wrapped ({"taskDefinition": {...}}) and unwrapped formats.
func LoadTaskDefinitionFromFile(filepath string) (*models.TaskDefinition, error) {
	src, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return ParseTaskDefinition(src)
}

// ParseTaskDefinition parses task definition JSON data.
// It attempts to parse as wrapped format first, then falls back to unwrapped format.
func ParseTaskDefinition(data []byte) (*models.TaskDefinition, error) {
	var taskDefinitionJson models.TaskDefinitionJson
	var taskDefinition models.TaskDefinition

	// Try parsing as wrapped format first
	if err := json.Unmarshal(data, &taskDefinitionJson); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	// If wrapped format has container definitions, use it
	if taskDefinitionJson.TaskDefinition.ContainerDefinitions != nil {
		return &taskDefinitionJson.TaskDefinition, nil
	}

	// Otherwise, try parsing as unwrapped format
	if err := json.Unmarshal(data, &taskDefinition); err != nil {
		return nil, fmt.Errorf("failed to parse as unwrapped JSON: %w", err)
	}

	return &taskDefinition, nil
}

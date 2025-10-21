package main

import (
	"encoding/json"
	"os"
	"strconv"
	"testing"
)

// TestParseWrappedJSON tests parsing of wrapped JSON format
func TestParseWrappedJSON(t *testing.T) {
	wrappedJSON := `{
		"taskDefinition": {
			"family": "test-family",
			"cpu": "512",
			"memory": "1024",
			"networkMode": "awsvpc",
			"requiresCompatibilities": ["FARGATE"],
			"executionRoleArn": "arn:aws:iam::123456789012:role/ecsTaskExecutionRole",
			"containerDefinitions": [
				{
					"name": "nginx",
					"image": "nginx:latest",
					"cpu": 256,
					"memory": 512,
					"memoryReservation": 256,
					"essential": true,
					"portMappings": [],
					"environment": [],
					"logConfiguration": {
						"logDriver": "awslogs"
					},
					"secrets": [],
					"command": [],
					"volumesFrom": []
				}
			],
			"placementConstraints": []
		}
	}`

	var taskDefinitionJson TaskDefinitionJson
	err := json.Unmarshal([]byte(wrappedJSON), &taskDefinitionJson)
	if err != nil {
		t.Fatalf("Failed to unmarshal wrapped JSON: %v", err)
	}

	if taskDefinitionJson.TaskDefinition.ContainerDefinitions == nil {
		t.Error("Expected ContainerDefinitions to be non-nil for wrapped JSON")
	}

	if taskDefinitionJson.TaskDefinition.Family != "test-family" {
		t.Errorf("Expected Family to be 'test-family', got '%s'", taskDefinitionJson.TaskDefinition.Family)
	}

	if taskDefinitionJson.TaskDefinition.CPU != "512" {
		t.Errorf("Expected CPU to be '512', got '%s'", taskDefinitionJson.TaskDefinition.CPU)
	}

	if taskDefinitionJson.TaskDefinition.Memory != "1024" {
		t.Errorf("Expected Memory to be '1024', got '%s'", taskDefinitionJson.TaskDefinition.Memory)
	}

	if len(taskDefinitionJson.TaskDefinition.ContainerDefinitions) != 1 {
		t.Errorf("Expected 1 container, got %d", len(taskDefinitionJson.TaskDefinition.ContainerDefinitions))
	}

	container := taskDefinitionJson.TaskDefinition.ContainerDefinitions[0]
	if container.Name != "nginx" {
		t.Errorf("Expected container name 'nginx', got '%s'", container.Name)
	}
	if container.CPU != 256 {
		t.Errorf("Expected container CPU 256, got %d", container.CPU)
	}
	if container.Memory != 512 {
		t.Errorf("Expected container Memory 512, got %d", container.Memory)
	}
	if container.MemoryReservation != 256 {
		t.Errorf("Expected container MemoryReservation 256, got %d", container.MemoryReservation)
	}
}

// TestParseUnwrappedJSON tests parsing of direct/unwrapped JSON format
func TestParseUnwrappedJSON(t *testing.T) {
	unwrappedJSON := `{
		"family": "test-family-unwrapped",
		"cpu": "256",
		"memory": "512",
		"networkMode": "bridge",
		"requiresCompatibilities": ["EC2"],
		"executionRoleArn": "arn:aws:iam::123456789012:role/ecsTaskExecutionRole",
		"containerDefinitions": [
			{
				"name": "app",
				"image": "app:latest",
				"cpu": 128,
				"memory": 256,
				"memoryReservation": 128,
				"essential": true,
				"portMappings": [],
				"environment": [],
				"logConfiguration": {
					"logDriver": "awslogs"
				},
				"secrets": [],
				"command": [],
				"volumesFrom": []
			}
		],
		"placementConstraints": []
	}`

	var taskDefinition TaskDefinition
	err := json.Unmarshal([]byte(unwrappedJSON), &taskDefinition)
	if err != nil {
		t.Fatalf("Failed to unmarshal unwrapped JSON: %v", err)
	}

	if taskDefinition.Family != "test-family-unwrapped" {
		t.Errorf("Expected Family to be 'test-family-unwrapped', got '%s'", taskDefinition.Family)
	}

	if taskDefinition.CPU != "256" {
		t.Errorf("Expected CPU to be '256', got '%s'", taskDefinition.CPU)
	}

	if taskDefinition.Memory != "512" {
		t.Errorf("Expected Memory to be '512', got '%s'", taskDefinition.Memory)
	}

	if len(taskDefinition.ContainerDefinitions) != 1 {
		t.Errorf("Expected 1 container, got %d", len(taskDefinition.ContainerDefinitions))
	}

	container := taskDefinition.ContainerDefinitions[0]
	if container.Name != "app" {
		t.Errorf("Expected container name 'app', got '%s'", container.Name)
	}
	if container.CPU != 128 {
		t.Errorf("Expected container CPU 128, got %d", container.CPU)
	}
}

// TestParseInvalidJSON tests error handling for invalid JSON
func TestParseInvalidJSON(t *testing.T) {
	invalidJSON := `{invalid json}`

	var taskDefinitionJson TaskDefinitionJson
	err := json.Unmarshal([]byte(invalidJSON), &taskDefinitionJson)
	if err == nil {
		t.Error("Expected error for invalid JSON, got nil")
	}
}

// TestParseEmptyJSON tests error handling for empty JSON
func TestParseEmptyJSON(t *testing.T) {
	emptyJSON := `{}`

	var taskDefinitionJson TaskDefinitionJson
	err := json.Unmarshal([]byte(emptyJSON), &taskDefinitionJson)
	if err != nil {
		t.Errorf("Unexpected error for empty JSON: %v", err)
	}

	// Empty JSON should result in nil ContainerDefinitions
	if taskDefinitionJson.TaskDefinition.ContainerDefinitions != nil {
		t.Error("Expected ContainerDefinitions to be nil for empty wrapped JSON")
	}
}

// TestResourceCalculation tests the calculation of total resources
func TestResourceCalculation(t *testing.T) {
	taskDefinition := TaskDefinition{
		CPU:    "1024",
		Memory: "2048",
		ContainerDefinitions: []ContainerDefinition{
			{
				Name:              "container1",
				CPU:               256,
				Memory:            512,
				MemoryReservation: 256,
			},
			{
				Name:              "container2",
				CPU:               512,
				Memory:            1024,
				MemoryReservation: 512,
			},
		},
	}

	var allCPU, allMemory, allMemoryReservation int64
	for _, container := range taskDefinition.ContainerDefinitions {
		allCPU += container.CPU
		allMemory += container.Memory
		allMemoryReservation += container.MemoryReservation
	}

	expectedCPU := int64(768)
	expectedMemory := int64(1536)
	expectedMemoryReservation := int64(768)

	if allCPU != expectedCPU {
		t.Errorf("Expected total CPU %d, got %d", expectedCPU, allCPU)
	}
	if allMemory != expectedMemory {
		t.Errorf("Expected total Memory %d, got %d", expectedMemory, allMemory)
	}
	if allMemoryReservation != expectedMemoryReservation {
		t.Errorf("Expected total MemoryReservation %d, got %d", expectedMemoryReservation, allMemoryReservation)
	}
}

// TestLeftoverCalculation tests the calculation of leftover resources
func TestLeftoverCalculation(t *testing.T) {
	tests := []struct {
		name                           string
		taskCPU                        int64
		taskMemory                     int64
		containerCPU                   int64
		containerMemory                int64
		containerMemoryReservation     int64
		expectedLeftoverCPU            int64
		expectedLeftoverMemory         int64
		expectedLeftoverMemReservation int64
	}{
		{
			name:                           "Positive leftover",
			taskCPU:                        1024,
			taskMemory:                     2048,
			containerCPU:                   512,
			containerMemory:                1024,
			containerMemoryReservation:     512,
			expectedLeftoverCPU:            512,
			expectedLeftoverMemory:         1024,
			expectedLeftoverMemReservation: 1536,
		},
		{
			name:                           "Negative leftover (over-allocation)",
			taskCPU:                        512,
			taskMemory:                     1024,
			containerCPU:                   768,
			containerMemory:                1536,
			containerMemoryReservation:     1280,
			expectedLeftoverCPU:            -256,
			expectedLeftoverMemory:         -512,
			expectedLeftoverMemReservation: -256,
		},
		{
			name:                           "Zero leftover (exact match)",
			taskCPU:                        1024,
			taskMemory:                     2048,
			containerCPU:                   1024,
			containerMemory:                2048,
			containerMemoryReservation:     2048,
			expectedLeftoverCPU:            0,
			expectedLeftoverMemory:         0,
			expectedLeftoverMemReservation: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			leftoverCPU := tt.taskCPU - tt.containerCPU
			leftoverMemory := tt.taskMemory - tt.containerMemory
			leftoverMemoryReservation := tt.taskMemory - tt.containerMemoryReservation

			if leftoverCPU != tt.expectedLeftoverCPU {
				t.Errorf("Expected leftover CPU %d, got %d", tt.expectedLeftoverCPU, leftoverCPU)
			}
			if leftoverMemory != tt.expectedLeftoverMemory {
				t.Errorf("Expected leftover Memory %d, got %d", tt.expectedLeftoverMemory, leftoverMemory)
			}
			if leftoverMemoryReservation != tt.expectedLeftoverMemReservation {
				t.Errorf("Expected leftover MemoryReservation %d, got %d", tt.expectedLeftoverMemReservation, leftoverMemoryReservation)
			}

			// Test negative value detection
			if leftoverCPU < 0 && tt.expectedLeftoverCPU >= 0 {
				t.Error("CPU leftover should be negative but is not")
			}
			if leftoverMemory < 0 && tt.expectedLeftoverMemory >= 0 {
				t.Error("Memory leftover should be negative but is not")
			}
		})
	}
}

// TestReadTestTaskFile tests reading the actual test_task.json file
func TestReadTestTaskFile(t *testing.T) {
	// Check if test_task.json exists
	if _, err := os.Stat("test_task.json"); os.IsNotExist(err) {
		t.Skip("test_task.json not found, skipping test")
	}

	src, err := os.ReadFile("test_task.json")
	if err != nil {
		t.Fatalf("Failed to read test_task.json: %v", err)
	}

	var taskDefinitionJson TaskDefinitionJson
	var taskDefinition TaskDefinition

	if err := json.Unmarshal(src, &taskDefinitionJson); err != nil {
		t.Fatalf("Failed to unmarshal test_task.json: %v", err)
	}

	if taskDefinitionJson.TaskDefinition.ContainerDefinitions == nil {
		if err = json.Unmarshal(src, &taskDefinition); err != nil {
			t.Fatalf("Failed to unmarshal as unwrapped JSON: %v", err)
		}
	} else {
		taskDefinition = taskDefinitionJson.TaskDefinition
	}

	// Verify basic structure
	if taskDefinition.ContainerDefinitions == nil {
		t.Error("ContainerDefinitions should not be nil")
	}

	if len(taskDefinition.ContainerDefinitions) == 0 {
		t.Error("ContainerDefinitions should not be empty")
	}

	// Verify that CPU and Memory are parseable
	if taskDefinition.CPU == "" {
		t.Error("Task CPU should not be empty")
	}
	if taskDefinition.Memory == "" {
		t.Error("Task Memory should not be empty")
	}
}

// TestNegativeLeftoverDetection tests detection of negative leftover values
func TestNegativeLeftoverDetection(t *testing.T) {
	testCases := []struct {
		leftover   int64
		isNegative bool
	}{
		{leftover: 100, isNegative: false},
		{leftover: 0, isNegative: false},
		{leftover: -1, isNegative: true},
		{leftover: -500, isNegative: true},
	}

	for _, tc := range testCases {
		result := tc.leftover < 0
		if result != tc.isNegative {
			t.Errorf("For leftover %d, expected isNegative=%v, got %v", tc.leftover, tc.isNegative, result)
		}
	}
}

// TestParseIntErrors tests error handling for invalid CPU/Memory values
func TestParseIntErrors(t *testing.T) {
	testCases := []struct {
		name     string
		value    string
		hasError bool
	}{
		{name: "Valid number", value: "1024", hasError: false},
		{name: "Invalid string", value: "invalid", hasError: true},
		{name: "Empty string", value: "", hasError: true},
		{name: "Float value", value: "10.5", hasError: true},
		{name: "Negative value", value: "-100", hasError: false}, // Negative is valid for parseInt
		{name: "Very large number", value: "999999999999", hasError: false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := strconv.ParseInt(tc.value, 10, 64)
			if tc.hasError && err == nil {
				t.Errorf("Expected error for value '%s', got nil", tc.value)
			}
			if !tc.hasError && err != nil {
				t.Errorf("Expected no error for value '%s', got %v", tc.value, err)
			}
		})
	}
}

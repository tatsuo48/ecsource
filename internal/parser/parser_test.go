package parser

import (
	"os"
	"testing"
)

func TestParseTaskDefinitionWrapped(t *testing.T) {
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

	taskDef, err := ParseTaskDefinition([]byte(wrappedJSON))
	if err != nil {
		t.Fatalf("Failed to parse wrapped JSON: %v", err)
	}

	if taskDef.Family != "test-family" {
		t.Errorf("Expected Family 'test-family', got '%s'", taskDef.Family)
	}

	if taskDef.CPU != "512" {
		t.Errorf("Expected CPU '512', got '%s'", taskDef.CPU)
	}

	if len(taskDef.ContainerDefinitions) != 1 {
		t.Errorf("Expected 1 container, got %d", len(taskDef.ContainerDefinitions))
	}
}

func TestParseTaskDefinitionUnwrapped(t *testing.T) {
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

	taskDef, err := ParseTaskDefinition([]byte(unwrappedJSON))
	if err != nil {
		t.Fatalf("Failed to parse unwrapped JSON: %v", err)
	}

	if taskDef.Family != "test-family-unwrapped" {
		t.Errorf("Expected Family 'test-family-unwrapped', got '%s'", taskDef.Family)
	}

	if taskDef.CPU != "256" {
		t.Errorf("Expected CPU '256', got '%s'", taskDef.CPU)
	}
}

func TestParseTaskDefinitionInvalid(t *testing.T) {
	invalidJSON := `{invalid json}`

	_, err := ParseTaskDefinition([]byte(invalidJSON))
	if err == nil {
		t.Error("Expected error for invalid JSON, got nil")
	}
}

func TestLoadTaskDefinitionFromFile(t *testing.T) {
	// Create a temporary test file
	tmpFile, err := os.CreateTemp("", "test_task_*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	testJSON := `{
		"taskDefinition": {
			"family": "temp-test",
			"cpu": "256",
			"memory": "512",
			"networkMode": "bridge",
			"containerDefinitions": [
				{
					"name": "test-container",
					"image": "test:latest",
					"cpu": 128,
					"memory": 256,
					"memoryReservation": 128,
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
		}
	}`

	if _, err := tmpFile.Write([]byte(testJSON)); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	taskDef, err := LoadTaskDefinitionFromFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to load task definition: %v", err)
	}

	if taskDef.Family != "temp-test" {
		t.Errorf("Expected Family 'temp-test', got '%s'", taskDef.Family)
	}
}

func TestLoadTaskDefinitionFromFileNotFound(t *testing.T) {
	_, err := LoadTaskDefinitionFromFile("nonexistent_file.json")
	if err == nil {
		t.Error("Expected error for nonexistent file, got nil")
	}
}

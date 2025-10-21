package models

// TaskDefinitionJson represents a wrapped task definition JSON structure
type TaskDefinitionJson struct {
	TaskDefinition TaskDefinition `json:"taskDefinition"`
}

// TaskDefinition represents an ECS task definition
type TaskDefinition struct {
	ContainerDefinitions    []ContainerDefinition `json:"containerDefinitions"`
	CPU                     string                `json:"cpu"`
	ExecutionRoleArn        string                `json:"executionRoleArn"`
	Family                  string                `json:"family"`
	Memory                  string                `json:"memory"`
	NetworkMode             string                `json:"networkMode"`
	PlacementConstraints    []interface{}         `json:"placementConstraints"`
	RequiresCompatibilities []string              `json:"requiresCompatibilities"`
	TaskRoleArn             string                `json:"taskRoleArn"`
}

// ContainerDefinition represents an ECS container definition
type ContainerDefinition struct {
	Image                  string                 `json:"image"`
	Name                   string                 `json:"name"`
	CPU                    int64                  `json:"cpu"`
	MemoryReservation      int64                  `json:"memoryReservation"`
	Memory                 int64                  `json:"memory"`
	Essential              bool                   `json:"essential"`
	Environment            []Environment          `json:"environment"`
	LogConfiguration       LogConfiguration       `json:"logConfiguration"`
	Secrets                []Secret               `json:"secrets"`
	Command                []string               `json:"command"`
	DockerLabels           *DockerLabels          `json:"dockerLabels,omitempty"`
	PortMappings           []PortMapping          `json:"portMappings"`
	ReadonlyRootFilesystem *bool                  `json:"readonlyRootFilesystem,omitempty"`
	StopTimeout            *int64                 `json:"stopTimeout,omitempty"`
	VolumesFrom            []interface{}          `json:"volumesFrom"`
	FirelensConfiguration  *FirelensConfiguration `json:"firelensConfiguration,omitempty"`
}

// DockerLabels represents Docker labels for a container
type DockerLabels struct {
	Environment string `json:"environment"`
	Service     string `json:"service"`
	App         string `json:"app"`
}

// Environment represents an environment variable
type Environment struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// FirelensConfiguration represents Firelens configuration
type FirelensConfiguration struct {
	Options FirelensConfigurationOptions `json:"options"`
	Type    string                       `json:"type"`
}

// FirelensConfigurationOptions represents Firelens configuration options
type FirelensConfigurationOptions struct {
	ConfigFileType  string `json:"config-file-type"`
	ConfigFileValue string `json:"config-file-value"`
}

// LogConfiguration represents logging configuration
type LogConfiguration struct {
	LogDriver string                   `json:"logDriver"`
	Options   *LogConfigurationOptions `json:"options,omitempty"`
}

// LogConfigurationOptions represents logging configuration options
type LogConfigurationOptions struct {
	AwslogsGroup        string `json:"awslogs-group"`
	AwslogsRegion       string `json:"awslogs-region"`
	AwslogsStreamPrefix string `json:"awslogs-stream-prefix"`
}

// PortMapping represents a port mapping configuration
type PortMapping struct {
	ContainerPort int64  `json:"containerPort"`
	HostPort      int64  `json:"hostPort"`
	Protocol      string `json:"protocol"`
}

// Secret represents a secret reference
type Secret struct {
	Name      string `json:"name"`
	ValueFrom string `json:"valueFrom"`
}

// ResourceSummary represents calculated resource totals and leftovers
type ResourceSummary struct {
	TaskCPU                   int64
	TaskMemory                int64
	TotalCPU                  int64
	TotalMemory               int64
	TotalMemoryReservation    int64
	LeftoverCPU               int64
	LeftoverMemory            int64
	LeftoverMemoryReservation int64
}

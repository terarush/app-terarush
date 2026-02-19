package request

type CreateNodeRequest struct {
	Name          string            `json:"name" validate:"required,min=3,max=255"`
	Image         string            `json:"image" validate:"required"`
	Port          int               `json:"port" validate:"required,min=1,max=65535"`
	InternalPort  int               `json:"internal_port" validate:"required,min=1,max=65535"`
	CPULimit      float64           `json:"cpu_limit" validate:"omitempty,min=0.1,max=32"`
	MemoryLimit   int64             `json:"memory_limit" validate:"omitempty,min=128,max=32768"` // in MB
	Environment   map[string]string `json:"environment"`
	Volumes       map[string]string `json:"volumes"`
	Command       string            `json:"command"`
	RestartPolicy string            `json:"restart_policy" validate:"omitempty,oneof=no on-failure always unless-stopped"`
}

type UpdateNodeRequest struct {
	Name          string            `json:"name" validate:"omitempty,min=3,max=255"`
	CPULimit      float64           `json:"cpu_limit" validate:"omitempty,min=0.1,max=32"`
	MemoryLimit   int64             `json:"memory_limit" validate:"omitempty,min=128,max=32768"`
	Environment   map[string]string `json:"environment"`
	Volumes       map[string]string `json:"volumes"`
	Command       string            `json:"command"`
	RestartPolicy string            `json:"restart_policy" validate:"omitempty,oneof=no on-failure always unless-stopped"`
}

type NodeActionRequest struct {
	Action string `json:"action" validate:"required,oneof=start stop restart"`
}

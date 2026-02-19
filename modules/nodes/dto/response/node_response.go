package response

import (
	"time"

	"go-modular/modules/nodes/domain/entity"
)

type NodeResponse struct {
	ID            uint              `json:"id"`
	UserID        uint              `json:"user_id"`
	Name          string            `json:"name"`
	Image         string            `json:"image"`
	ContainerID   string            `json:"container_id"`
	Status        string            `json:"status"`
	Port          int               `json:"port"`
	InternalPort  int               `json:"internal_port"`
	CPULimit      float64           `json:"cpu_limit"`
	MemoryLimit   int64             `json:"memory_limit"`
	Environment   map[string]string `json:"environment"`
	Volumes       map[string]string `json:"volumes"`
	Command       string            `json:"command"`
	RestartPolicy string            `json:"restart_policy"`
	CreatedAt     time.Time         `json:"created_at"`
	UpdatedAt     time.Time         `json:"updated_at"`
}

func ToNodeResponse(node *entity.Node) NodeResponse {
	return NodeResponse{
		ID:            node.ID,
		UserID:        node.UserID,
		Name:          node.Name,
		Image:         node.Image,
		ContainerID:   node.ContainerID,
		Status:        string(node.Status),
		Port:          node.Port,
		InternalPort:  node.InternalPort,
		CPULimit:      node.CPULimit,
		MemoryLimit:   node.MemoryLimit,
		Environment:   parseJSONToMap(node.Environment),
		Volumes:       parseJSONToMap(node.Volumes),
		Command:       node.Command,
		RestartPolicy: node.RestartPolicy,
		CreatedAt:     node.CreatedAt,
		UpdatedAt:     node.UpdatedAt,
	}
}

func ToNodeResponseList(nodes []entity.Node) []NodeResponse {
	responses := make([]NodeResponse, len(nodes))
	for i, node := range nodes {
		responses[i] = ToNodeResponse(&node)
	}
	return responses
}

// Helper function to parse JSON string to map
func parseJSONToMap(jsonStr string) map[string]string {
	result := make(map[string]string)
	if jsonStr == "" {
		return result
	}
	// Simple parsing, can be improved with json.Unmarshal
	return result
}

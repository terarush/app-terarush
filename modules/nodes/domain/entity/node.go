package entity

import (
	"time"

	"gorm.io/gorm"
)

type NodeStatus string

const (
	NodeStatusRunning NodeStatus = "running"
	NodeStatusStopped NodeStatus = "stopped"
	NodeStatusCreated NodeStatus = "created"
	NodeStatusError   NodeStatus = "error"
)

type Node struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	UserID        uint           `gorm:"not null;index" json:"user_id"`
	Name          string         `gorm:"not null;size:255" json:"name"`
	Image         string         `gorm:"not null;size:255" json:"image"`
	ContainerID   string         `gorm:"size:255;uniqueIndex" json:"container_id"`
	Status        NodeStatus     `gorm:"type:varchar(20);default:'created'" json:"status"`
	Port          int            `gorm:"not null" json:"port"`
	InternalPort  int            `gorm:"not null" json:"internal_port"`
	CPULimit      float64        `gorm:"default:1.0" json:"cpu_limit"`
	MemoryLimit   int64          `gorm:"default:512" json:"memory_limit"` // in MB
	Environment   string         `gorm:"type:text" json:"environment"`    // JSON string
	Volumes       string         `gorm:"type:text" json:"volumes"`        // JSON string
	Command       string         `gorm:"type:text" json:"command"`
	RestartPolicy string         `gorm:"size:50;default:'unless-stopped'" json:"restart_policy"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (Node) TableName() string {
	return "nodes"
}

package entity

import (
	"time"
)

// Product represents a container product
type Product struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	Name        string     `gorm:"size:255;not null" json:"name"`
	Description string     `gorm:"type:text" json:"description"`
	Price       float64    `gorm:"type:decimal(15,2);not null" json:"price"`
	CPUCores    int        `gorm:"not null;default:1" json:"cpu_cores"`
	RAMMB       int        `gorm:"not null;default:512" json:"ram_mb"`
	StorageGB   int        `gorm:"not null;default:10" json:"storage_gb"`
	BandwidthGB int        `gorm:"not null;default:100" json:"bandwidth_gb"`
	IsActive    bool       `gorm:"default:true" json:"is_active"`
	Stock       int        `gorm:"default:-1" json:"stock"` // -1 means unlimited
	ImageURL    string     `gorm:"size:255" json:"image_url"`
	DockerImage string     `gorm:"size:100" json:"docker_image"` // e.g., "node", "python", "ubuntu"
	ImageTags   string     `gorm:"type:json" json:"image_tags"`  // JSON array: ["18","20","latest"]
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `gorm:"index" json:"deleted_at,omitempty"`
}

// TableName specifies the table name for Product
func (Product) TableName() string {
	return "products"
}

// IsAvailable checks if the product is available for purchase
func (p *Product) IsAvailable() bool {
	return p.IsActive && (p.Stock == -1 || p.Stock > 0)
}

// DecrementStock reduces stock by quantity
func (p *Product) DecrementStock(quantity int) error {
	if p.Stock == -1 {
		return nil // Unlimited stock
	}

	if p.Stock < quantity {
		return ErrInsufficientStock
	}

	p.Stock -= quantity
	return nil
}

// Custom errors
var (
	ErrInsufficientStock = &CustomError{Message: "insufficient stock"}
)

type CustomError struct {
	Message string
}

func (e *CustomError) Error() string {
	return e.Message
}

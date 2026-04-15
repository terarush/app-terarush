package response

import (
	"encoding/json"
	"go-modular/modules/products/domain/entity"
	"time"
)

// ProductResponse represents the product response
type ProductResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	CPUCores    int       `json:"cpu_cores"`
	RAMMB       int       `json:"ram_mb"`
	StorageGB   int       `json:"storage_gb"`
	BandwidthGB int       `json:"bandwidth_gb"`
	IsActive    bool      `json:"is_active"`
	Stock       int       `json:"stock"`
	ImageURL    string    `json:"image_url"`
	DockerImage string    `json:"docker_image"`
	ImageTags   []string  `json:"image_tags"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// FromEntity converts entity to response
func FromEntity(product *entity.Product) *ProductResponse {
	// Parse image_tags JSON string to []string
	var imageTags []string
	if product.ImageTags != "" {
		json.Unmarshal([]byte(product.ImageTags), &imageTags)
	}

	return &ProductResponse{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		CPUCores:    product.CPUCores,
		RAMMB:       product.RAMMB,
		StorageGB:   product.StorageGB,
		BandwidthGB: product.BandwidthGB,
		IsActive:    product.IsActive,
		Stock:       product.Stock,
		ImageURL:    product.ImageURL,
		DockerImage: product.DockerImage,
		ImageTags:   imageTags,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
	}
}

// FromEntities converts multiple entities to responses
func FromEntities(products []*entity.Product) []*ProductResponse {
	responses := make([]*ProductResponse, len(products))
	for i, product := range products {
		responses[i] = FromEntity(product)
	}
	return responses
}

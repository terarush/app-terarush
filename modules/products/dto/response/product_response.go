package response

import (
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
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// FromEntity converts entity to response
func FromEntity(product *entity.Product) *ProductResponse {
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

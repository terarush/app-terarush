package request

// CreateProductRequest represents the request to create a product
type CreateProductRequest struct {
	Name        string  `json:"name" validate:"required,min=3,max=255"`
	Description string  `json:"description"`
	Price       float64 `json:"price" validate:"required,gt=0"`
	CPUCores    int     `json:"cpu_cores" validate:"required,gte=1"`
	RAMMB       int     `json:"ram_mb" validate:"required,gte=128"`
	StorageGB   int     `json:"storage_gb" validate:"required,gte=1"`
	BandwidthGB int     `json:"bandwidth_gb" validate:"required,gte=1"`
	IsActive    bool    `json:"is_active"`
	Stock       int     `json:"stock" validate:"gte=-1"` // -1 means unlimited
	ImageURL    string  `json:"image_url"`
}

// UpdateProductRequest represents the request to update a product
type UpdateProductRequest struct {
	Name        string  `json:"name" validate:"required,min=3,max=255"`
	Description string  `json:"description"`
	Price       float64 `json:"price" validate:"required,gt=0"`
	CPUCores    int     `json:"cpu_cores" validate:"required,gte=1"`
	RAMMB       int     `json:"ram_mb" validate:"required,gte=128"`
	StorageGB   int     `json:"storage_gb" validate:"required,gte=1"`
	BandwidthGB int     `json:"bandwidth_gb" validate:"required,gte=1"`
	IsActive    bool    `json:"is_active"`
	Stock       int     `json:"stock" validate:"gte=-1"`
	ImageURL    string  `json:"image_url"`
}

package response

import (
	"go-modular/modules/assets/domain/entity"
	"time"
)

// AssetResponse is the response for asset data
type AssetResponse struct {
	ID        uint      `json:"id"`
	URL       string    `json:"url"`
	Path      string    `json:"path"`
	FileName  string    `json:"file_name"`
	MimeType  string    `json:"mime_type"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// FromEntity converts asset entity to response
func FromEntity(asset *entity.Asset) AssetResponse {
	return AssetResponse{
		ID:        asset.ID,
		URL:       asset.URL,
		Path:      asset.URL,
		FileName:  asset.FileName,
		MimeType:  asset.MimeType,
		CreatedAt: asset.CreatedAt,
		UpdatedAt: asset.UpdatedAt,
	}
}

// FromEntities converts multiple asset entities to responses
func FromEntities(assets []*entity.Asset) []AssetResponse {
	responses := make([]AssetResponse, len(assets))
	for i, asset := range assets {
		responses[i] = FromEntity(asset)
	}
	return responses
}

// UploadAssetResponse is the response for asset upload
type UploadAssetResponse struct {
	ID   uint   `json:"id"`
	URL  string `json:"url"`
	Path string `json:"path"`
}

// ListAssetsResponse is the response for listing assets
type ListAssetsResponse struct {
	Assets   []AssetResponse `json:"assets"`
	Total    int64           `json:"total"`
	Page     int             `json:"page"`
	PageSize int             `json:"page_size"`
}

package response

import (
	"go-modular/modules/assets/domain/entity"
	"time"
)

// AssetResponse is the response for asset data
type AssetResponse struct {
	ID        uint      `json:"id"`
	URL       string    `json:"url"`
	CreatedAt time.Time `json:"created_at"`
}

// FromEntity converts asset entity to response
func FromEntity(asset *entity.Asset) AssetResponse {
	return AssetResponse{
		ID:        asset.ID,
		URL:       asset.URL,
		CreatedAt: asset.CreatedAt,
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

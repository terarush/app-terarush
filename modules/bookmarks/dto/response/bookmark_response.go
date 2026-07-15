package response

import (
	"go-modular/modules/bookmarks/domain/entity"
	"time"
)

type BookmarkResponse struct {
	ID           uint      `json:"id"`
	ResourceType string    `json:"resource_type"`
	ResourceID   uint      `json:"resource_id"`
	CollectionID *uint     `json:"collection_id,omitempty"`
	Note         string    `json:"note,omitempty"`
	Tags         string    `json:"tags,omitempty"`
	Progress     float64   `json:"progress"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
}

func BmFromEntity(b *entity.Bookmark) *BookmarkResponse {
	return &BookmarkResponse{
		ID: b.ID, ResourceType: b.ResourceType, ResourceID: b.ResourceID,
		CollectionID: b.CollectionID, Note: b.Note, Tags: b.Tags,
		Progress: b.Progress, Status: b.Status, CreatedAt: b.CreatedAt,
	}
}

func BmFromEntities(items []*entity.Bookmark) []*BookmarkResponse {
	res := make([]*BookmarkResponse, len(items))
	for i, v := range items { res[i] = BmFromEntity(v) }
	return res
}

type CollectionResponse struct {
	ID     uint   `json:"id"`
	Name   string `json:"name"`
	Icon   string `json:"icon,omitempty"`
	Public bool   `json:"is_public"`
}

type GoalResponse struct {
	Year     int   `json:"year"`
	Goal     int   `json:"goal"`
	Progress int   `json:"progress"`
}

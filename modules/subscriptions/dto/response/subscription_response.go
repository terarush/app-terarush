package response

import (
	"go-modular/modules/subscriptions/domain/entity"
	"time"
)

type SubscriptionResponse struct {
	ID           uint      `json:"id"`
	ResourceType string    `json:"resource_type"`
	ResourceID   uint      `json:"resource_id"`
	Status       string    `json:"status"`
	Frequency    string    `json:"frequency"`
	CreatedAt    time.Time `json:"created_at"`
}

func SubFromEntity(s *entity.Subscription) *SubscriptionResponse {
	return &SubscriptionResponse{
		ID: s.ID, ResourceType: s.ResourceType, ResourceID: s.ResourceID,
		Status: s.Status, Frequency: s.Frequency, CreatedAt: s.CreatedAt,
	}
}

func SubFromEntities(items []*entity.Subscription) []*SubscriptionResponse {
	res := make([]*SubscriptionResponse, len(items))
	for i, v := range items { res[i] = SubFromEntity(v) }
	return res
}

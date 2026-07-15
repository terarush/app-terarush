package response

import (
	"go-modular/modules/notifications/domain/entity"
	"time"
)

type NotificationResponse struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"user_id"`
	Type      string    `json:"type"`
	Title     string    `json:"title"`
	Body      string    `json:"body,omitempty"`
	Link      string    `json:"link,omitempty"`
	Icon      string    `json:"icon,omitempty"`
	IsRead    bool      `json:"is_read"`
	Priority  string    `json:"priority"`
	CreatedAt time.Time `json:"created_at"`
	ReadAt    *time.Time `json:"read_at,omitempty"`
}

func NotifFromEntity(n *entity.Notification) *NotificationResponse {
	return &NotificationResponse{
		ID: n.ID, UserID: n.UserID, Type: n.Type, Title: n.Title,
		Body: n.Body, Link: n.Link, Icon: n.Icon, IsRead: n.IsRead,
		Priority: n.Priority, CreatedAt: n.CreatedAt, ReadAt: n.ReadAt,
	}
}

func NotifFromEntities(items []*entity.Notification) []*NotificationResponse {
	res := make([]*NotificationResponse, len(items))
	for i, v := range items { res[i] = NotifFromEntity(v) }
	return res
}

type UnreadCountResponse struct {
	Count int64 `json:"count"`
}

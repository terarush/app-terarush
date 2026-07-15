package repository

import (
	"context"
	"go-modular/modules/notifications/domain/entity"
)

type NotificationRepository interface {
	FindAll(ctx context.Context, filter map[string]any) ([]*entity.Notification, error)
	FindByID(ctx context.Context, id uint) (*entity.Notification, error)
	FindByUserID(ctx context.Context, userID uint, status string, limit int) ([]*entity.Notification, error)
	CountUnread(ctx context.Context, userID uint) (int64, error)
	MarkRead(ctx context.Context, id uint) error
	MarkAllRead(ctx context.Context, userID uint) error
	Create(ctx context.Context, n *entity.Notification) error
	Update(ctx context.Context, n *entity.Notification) error
	Delete(ctx context.Context, id uint) error
	DeleteOlderThan(ctx context.Context, days int) (int64, error)
}

type NotificationTemplateRepository interface {
	FindAll(ctx context.Context) ([]*entity.NotificationTemplate, error)
	FindByID(ctx context.Context, id uint) (*entity.NotificationTemplate, error)
	FindByName(ctx context.Context, name string) (*entity.NotificationTemplate, error)
	Create(ctx context.Context, t *entity.NotificationTemplate) error
	Update(ctx context.Context, t *entity.NotificationTemplate) error
	Delete(ctx context.Context, id uint) error
}

type NotificationPreferenceRepository interface {
	FindByUserID(ctx context.Context, userID uint) (*entity.NotificationPreference, error)
	Upsert(ctx context.Context, p *entity.NotificationPreference) error
}

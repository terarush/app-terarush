package repository

import (
	"context"
	"go-modular/modules/analytics/domain/entity"
)

type AnalyticsEventRepository interface {
	FindAll(ctx context.Context, filter map[string]any) ([]*entity.AnalyticsEvent, error)
	FindByID(ctx context.Context, id uint) (*entity.AnalyticsEvent, error)
	Create(ctx context.Context, e *entity.AnalyticsEvent) error
	DeleteOlderThan(ctx context.Context, days int) (int64, error)
	FindByEventType(ctx context.Context, eventType string, since string, limit int) ([]*entity.AnalyticsEvent, error)
	CountByEventType(ctx context.Context, eventType string, since string) (int64, error)
}

type BlogStatsRepository interface {
	FindByBlogID(ctx context.Context, blogID uint) (*entity.BlogStats, error)
	Upsert(ctx context.Context, s *entity.BlogStats) error
	IncrementView(ctx context.Context, blogID uint) error
	IncrementUniqueView(ctx context.Context, blogID uint) error
	TopBlogs(ctx context.Context, limit int) ([]*entity.BlogStats, error)
}

type UserAnalyticsRepository interface {
	FindByUserID(ctx context.Context, userID uint) (*entity.UserAnalytics, error)
	Upsert(ctx context.Context, u *entity.UserAnalytics) error
	UpdateLastActive(ctx context.Context, userID uint) error
}

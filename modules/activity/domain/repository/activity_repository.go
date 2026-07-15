package repository

import (
	"context"
	"go-modular/modules/activity/domain/entity"
)

type ActivityRepository interface {
	FindAll(ctx context.Context, filter map[string]any) ([]*entity.Activity, error)
	FindByID(ctx context.Context, id uint) (*entity.Activity, error)
	FindByUserID(ctx context.Context, userID uint, limit int) ([]*entity.Activity, error)
	Create(ctx context.Context, a *entity.Activity) error
	Delete(ctx context.Context, id uint) error
	DeleteOlderThan(ctx context.Context, days int) (int64, error)
}

type ActivityFeedRepository interface {
	FindByUserID(ctx context.Context, userID uint, limit int) ([]*entity.ActivityFeed, error)
	Create(ctx context.Context, f *entity.ActivityFeed) error
	MarkRead(ctx context.Context, id uint) error
	MarkAllRead(ctx context.Context, userID uint) error
	CountUnread(ctx context.Context, userID uint) (int64, error)
	DeleteOlderThan(ctx context.Context, days int) (int64, error)
}

type FollowRepository interface {
	FindAll(ctx context.Context, filter map[string]any) ([]*entity.Follow, error)
	FindByFollowerID(ctx context.Context, followerID uint) ([]*entity.Follow, error)
	FindByFollowingID(ctx context.Context, followingID uint) ([]*entity.Follow, error)
	Create(ctx context.Context, f *entity.Follow) error
	Delete(ctx context.Context, followerID, followingID uint) error
	IsFollowing(ctx context.Context, followerID, followingID uint) (bool, error)
	CountFollowers(ctx context.Context, userID uint) (int64, error)
	CountFollowing(ctx context.Context, userID uint) (int64, error)
}

type EngagementMetricRepository interface {
	FindByUserID(ctx context.Context, userID uint) (*entity.EngagementMetric, error)
	Upsert(ctx context.Context, m *entity.EngagementMetric) error
	TopUsers(ctx context.Context, limit int) ([]*entity.EngagementMetric, error)
}

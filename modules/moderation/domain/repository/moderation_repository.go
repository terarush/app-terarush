package repository

import (
	"context"
	"go-modular/modules/moderation/domain/entity"
)

type ReportRepository interface {
	FindAll(ctx context.Context, filter map[string]any) ([]*entity.Report, error)
	FindByID(ctx context.Context, id uint) (*entity.Report, error)
	FindByStatus(ctx context.Context, status string, limit int) ([]*entity.Report, error)
	Create(ctx context.Context, r *entity.Report) error
	Update(ctx context.Context, r *entity.Report) error
	Delete(ctx context.Context, id uint) error
	CountByStatus(ctx context.Context, status string) (int64, error)
}

type ModerationActionRepository interface {
	FindAll(ctx context.Context, filter map[string]any) ([]*entity.ModerationAction, error)
	FindByUserID(ctx context.Context, userID uint) ([]*entity.ModerationAction, error)
	Create(ctx context.Context, a *entity.ModerationAction) error
}

type BannedUserRepository interface {
	FindByUserID(ctx context.Context, userID uint) (*entity.BannedUser, error)
	FindAllActive(ctx context.Context) ([]*entity.BannedUser, error)
	Create(ctx context.Context, b *entity.BannedUser) error
	Delete(ctx context.Context, userID uint) error
	RemoveExpired(ctx context.Context) (int64, error)
}

type ContentFilterRepository interface {
	FindAll(ctx context.Context, activeOnly bool) ([]*entity.ContentFilter, error)
	FindByType(ctx context.Context, filterType string) ([]*entity.ContentFilter, error)
	FindByID(ctx context.Context, id uint) (*entity.ContentFilter, error)
	Create(ctx context.Context, f *entity.ContentFilter) error
	Update(ctx context.Context, f *entity.ContentFilter) error
	Delete(ctx context.Context, id uint) error
	IncrementHit(ctx context.Context, id uint) error
}

package repository

import (
	"context"
	"errors"
	"go-modular/internal/pkg/database"
	"gorm.io/gorm"
	"go-modular/modules/analytics/domain/entity"
	"time"
)

var (
	ErrEventNotFound   = errors.New("event not found")
	ErrBlogStatsNotFound = errors.New("blog stats not found")
)

type AnalyticsEventRepositoryImpl struct{}

func (r AnalyticsEventRepositoryImpl) FindAll(ctx context.Context, filter map[string]any) ([]*entity.AnalyticsEvent, error) {
	var items []*entity.AnalyticsEvent
	query := database.DB.WithContext(ctx)
	for k, v := range filter {
		query = query.Where(k, v)
	}
	result := query.Order("created_at DESC").Find(&items)
	return items, result.Error
}

func (r AnalyticsEventRepositoryImpl) FindByID(ctx context.Context, id uint) (*entity.AnalyticsEvent, error) {
	var e entity.AnalyticsEvent
	err := database.DB.WithContext(ctx).First(&e, id).Error
	if err != nil {
		return nil, ErrEventNotFound
	}
	return &e, nil
}

func (r AnalyticsEventRepositoryImpl) Create(ctx context.Context, e *entity.AnalyticsEvent) error {
	return database.DB.WithContext(ctx).Create(e).Error
}

func (r AnalyticsEventRepositoryImpl) DeleteOlderThan(ctx context.Context, days int) (int64, error) {
	cutoff := time.Now().AddDate(0, 0, -days)
	result := database.DB.WithContext(ctx).Where("created_at < ?", cutoff).Delete(&entity.AnalyticsEvent{})
	return result.RowsAffected, result.Error
}

func (r AnalyticsEventRepositoryImpl) FindByEventType(ctx context.Context, eventType string, since string, limit int) ([]*entity.AnalyticsEvent, error) {
	var items []*entity.AnalyticsEvent
	query := database.DB.WithContext(ctx).Where("event_type = ?", eventType)
	if since != "" {
		query = query.Where("created_at >= ?", since)
	}
	if limit > 0 {
		query = query.Limit(limit)
	}
	result := query.Order("created_at DESC").Find(&items)
	return items, result.Error
}

func (r AnalyticsEventRepositoryImpl) CountByEventType(ctx context.Context, eventType string, since string) (int64, error) {
	var count int64
	query := database.DB.WithContext(ctx).Model(&entity.AnalyticsEvent{}).Where("event_type = ?", eventType)
	if since != "" {
		query = query.Where("created_at >= ?", since)
	}
	err := query.Count(&count).Error
	return count, err
}

func NewAnalyticsEventRepository() AnalyticsEventRepository {
	return AnalyticsEventRepositoryImpl{}
}

type BlogStatsRepositoryImpl struct{}

func (r BlogStatsRepositoryImpl) FindByBlogID(ctx context.Context, blogID uint) (*entity.BlogStats, error) {
	var s entity.BlogStats
	err := database.DB.WithContext(ctx).Where("blog_id = ?", blogID).First(&s).Error
	if err != nil {
		return nil, ErrBlogStatsNotFound
	}
	return &s, nil
}

func (r BlogStatsRepositoryImpl) Upsert(ctx context.Context, s *entity.BlogStats) error {
	return database.DB.WithContext(ctx).Save(s).Error
}

func (r BlogStatsRepositoryImpl) IncrementView(ctx context.Context, blogID uint) error {
	return database.DB.WithContext(ctx).Model(&entity.BlogStats{}).
		Where("blog_id = ?", blogID).
		UpdateColumn("view_count", gorm.Expr("view_count + 1")).Error
}

func (r BlogStatsRepositoryImpl) IncrementUniqueView(ctx context.Context, blogID uint) error {
	return database.DB.WithContext(ctx).Model(&entity.BlogStats{}).
		Where("blog_id = ?", blogID).
		UpdateColumn("unique_views", gorm.Expr("unique_views + 1")).Error
}

func (r BlogStatsRepositoryImpl) TopBlogs(ctx context.Context, limit int) ([]*entity.BlogStats, error) {
	var items []*entity.BlogStats
	err := database.DB.WithContext(ctx).Order("view_count DESC").Limit(limit).Find(&items).Error
	return items, err
}

func NewBlogStatsRepository() BlogStatsRepository {
	return BlogStatsRepositoryImpl{}
}

type UserAnalyticsRepositoryImpl struct{}

func (r UserAnalyticsRepositoryImpl) FindByUserID(ctx context.Context, userID uint) (*entity.UserAnalytics, error) {
	var u entity.UserAnalytics
	err := database.DB.WithContext(ctx).Where("user_id = ?", userID).First(&u).Error
	return &u, err
}

func (r UserAnalyticsRepositoryImpl) Upsert(ctx context.Context, u *entity.UserAnalytics) error {
	return database.DB.WithContext(ctx).Save(u).Error
}

func (r UserAnalyticsRepositoryImpl) UpdateLastActive(ctx context.Context, userID uint) error {
	now := time.Now()
	return database.DB.WithContext(ctx).Model(&entity.UserAnalytics{}).
		Where("user_id = ?", userID).
		UpdateColumn("last_active_at", now).Error
}

func NewUserAnalyticsRepository() UserAnalyticsRepository {
	return UserAnalyticsRepositoryImpl{}
}

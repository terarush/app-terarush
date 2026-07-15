package repository

import (
	"context"
	"errors"
	"go-modular/internal/pkg/database"
	"go-modular/modules/activity/domain/entity"
	"time"
)

var (
	ErrActivityNotFound   = errors.New("activity not found")
	ErrFeedNotFound       = errors.New("feed entry not found")
	ErrFollowNotFound     = errors.New("follow relationship not found")
	ErrMetricNotFound     = errors.New("engagement metric not found")
)

type ActivityRepositoryImpl struct{}

func (r ActivityRepositoryImpl) FindAll(ctx context.Context, filter map[string]any) ([]*entity.Activity, error) {
	var items []*entity.Activity
	query := database.DB.WithContext(ctx)
	for k, v := range filter {
		query = query.Where(k, v)
	}
	return items, query.Order("created_at DESC").Find(&items).Error
}

func (r ActivityRepositoryImpl) FindByID(ctx context.Context, id uint) (*entity.Activity, error) {
	var a entity.Activity
	err := database.DB.WithContext(ctx).First(&a, id).Error
	if err != nil {
		return nil, ErrActivityNotFound
	}
	return &a, nil
}

func (r ActivityRepositoryImpl) FindByUserID(ctx context.Context, userID uint, limit int) ([]*entity.Activity, error) {
	var items []*entity.Activity
	query := database.DB.WithContext(ctx).Where("actor_id = ? OR user_id = ?", userID, userID)
	if limit > 0 {
		query = query.Limit(limit)
	}
	return items, query.Order("created_at DESC").Find(&items).Error
}

func (r ActivityRepositoryImpl) Create(ctx context.Context, a *entity.Activity) error {
	return database.DB.WithContext(ctx).Create(a).Error
}

func (r ActivityRepositoryImpl) Delete(ctx context.Context, id uint) error {
	return database.DB.WithContext(ctx).Delete(&entity.Activity{}, id).Error
}

func (r ActivityRepositoryImpl) DeleteOlderThan(ctx context.Context, days int) (int64, error) {
	cutoff := time.Now().AddDate(0, 0, -days)
	result := database.DB.WithContext(ctx).Where("created_at < ?", cutoff).Delete(&entity.Activity{})
	return result.RowsAffected, result.Error
}

func NewActivityRepository() ActivityRepository {
	return ActivityRepositoryImpl{}
}

type ActivityFeedRepositoryImpl struct{}

func (r ActivityFeedRepositoryImpl) FindByUserID(ctx context.Context, userID uint, limit int) ([]*entity.ActivityFeed, error) {
	var items []*entity.ActivityFeed
	query := database.DB.WithContext(ctx).Where("user_id = ?", userID)
	if limit > 0 {
		query = query.Limit(limit)
	}
	return items, query.Order("created_at DESC").Find(&items).Error
}

func (r ActivityFeedRepositoryImpl) Create(ctx context.Context, f *entity.ActivityFeed) error {
	return database.DB.WithContext(ctx).Create(f).Error
}

func (r ActivityFeedRepositoryImpl) MarkRead(ctx context.Context, id uint) error {
	return database.DB.WithContext(ctx).Model(&entity.ActivityFeed{}).Where("id = ?", id).Update("is_read", true).Error
}

func (r ActivityFeedRepositoryImpl) MarkAllRead(ctx context.Context, userID uint) error {
	return database.DB.WithContext(ctx).Model(&entity.ActivityFeed{}).Where("user_id = ? AND is_read = ?", userID, false).Update("is_read", true).Error
}

func (r ActivityFeedRepositoryImpl) CountUnread(ctx context.Context, userID uint) (int64, error) {
	var count int64
	err := database.DB.WithContext(ctx).Model(&entity.ActivityFeed{}).Where("user_id = ? AND is_read = ?", userID, false).Count(&count).Error
	return count, err
}

func (r ActivityFeedRepositoryImpl) DeleteOlderThan(ctx context.Context, days int) (int64, error) {
	cutoff := time.Now().AddDate(0, 0, -days)
	result := database.DB.WithContext(ctx).Where("created_at < ?", cutoff).Delete(&entity.ActivityFeed{})
	return result.RowsAffected, result.Error
}

func NewActivityFeedRepository() ActivityFeedRepository {
	return ActivityFeedRepositoryImpl{}
}

type FollowRepositoryImpl struct{}

func (r FollowRepositoryImpl) FindAll(ctx context.Context, filter map[string]any) ([]*entity.Follow, error) {
	var items []*entity.Follow
	query := database.DB.WithContext(ctx)
	for k, v := range filter {
		query = query.Where(k, v)
	}
	return items, query.Order("created_at DESC").Find(&items).Error
}

func (r FollowRepositoryImpl) FindByFollowerID(ctx context.Context, followerID uint) ([]*entity.Follow, error) {
	var items []*entity.Follow
	err := database.DB.WithContext(ctx).Where("follower_id = ?", followerID).Find(&items).Error
	return items, err
}

func (r FollowRepositoryImpl) FindByFollowingID(ctx context.Context, followingID uint) ([]*entity.Follow, error) {
	var items []*entity.Follow
	err := database.DB.WithContext(ctx).Where("following_id = ?", followingID).Find(&items).Error
	return items, err
}

func (r FollowRepositoryImpl) Create(ctx context.Context, f *entity.Follow) error {
	return database.DB.WithContext(ctx).Create(f).Error
}

func (r FollowRepositoryImpl) Delete(ctx context.Context, followerID, followingID uint) error {
	return database.DB.WithContext(ctx).
		Where("follower_id = ? AND following_id = ?", followerID, followingID).
		Delete(&entity.Follow{}).Error
}

func (r FollowRepositoryImpl) IsFollowing(ctx context.Context, followerID, followingID uint) (bool, error) {
	var count int64
	err := database.DB.WithContext(ctx).Model(&entity.Follow{}).
		Where("follower_id = ? AND following_id = ?", followerID, followingID).
		Count(&count).Error
	return count > 0, err
}

func (r FollowRepositoryImpl) CountFollowers(ctx context.Context, userID uint) (int64, error) {
	var count int64
	err := database.DB.WithContext(ctx).Model(&entity.Follow{}).Where("following_id = ?", userID).Count(&count).Error
	return count, err
}

func (r FollowRepositoryImpl) CountFollowing(ctx context.Context, userID uint) (int64, error) {
	var count int64
	err := database.DB.WithContext(ctx).Model(&entity.Follow{}).Where("follower_id = ?", userID).Count(&count).Error
	return count, err
}

func NewFollowRepository() FollowRepository {
	return FollowRepositoryImpl{}
}

type EngagementMetricRepositoryImpl struct{}

func (r EngagementMetricRepositoryImpl) FindByUserID(ctx context.Context, userID uint) (*entity.EngagementMetric, error) {
	var m entity.EngagementMetric
	err := database.DB.WithContext(ctx).Where("user_id = ?", userID).First(&m).Error
	if err != nil {
		return nil, ErrMetricNotFound
	}
	return &m, nil
}

func (r EngagementMetricRepositoryImpl) Upsert(ctx context.Context, m *entity.EngagementMetric) error {
	return database.DB.WithContext(ctx).Save(m).Error
}

func (r EngagementMetricRepositoryImpl) TopUsers(ctx context.Context, limit int) ([]*entity.EngagementMetric, error) {
	var items []*entity.EngagementMetric
	err := database.DB.WithContext(ctx).Order("xp DESC").Limit(limit).Find(&items).Error
	return items, err
}

func NewEngagementMetricRepository() EngagementMetricRepository {
	return EngagementMetricRepositoryImpl{}
}

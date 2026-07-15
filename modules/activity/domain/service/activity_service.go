package service

import (
	"context"
	"errors"
	"go-modular/modules/activity/domain/entity"
	"go-modular/modules/activity/domain/repository"
	"time"
)

var (
	ErrActivityNotFound = errors.New("activity not found")
)

type ActivityService struct {
	actRepo  repository.ActivityRepository
	feedRepo repository.ActivityFeedRepository
	followRepo repository.FollowRepository
	metricRepo repository.EngagementMetricRepository
}

func NewActivityService(actRepo repository.ActivityRepository, feedRepo repository.ActivityFeedRepository, followRepo repository.FollowRepository, metricRepo repository.EngagementMetricRepository) *ActivityService {
	return &ActivityService{actRepo: actRepo, feedRepo: feedRepo, followRepo: followRepo, metricRepo: metricRepo}
}

// Activities
func (s *ActivityService) List(ctx context.Context, filter map[string]any) ([]*entity.Activity, error) {
	return s.actRepo.FindAll(ctx, filter)
}

func (s *ActivityService) Get(ctx context.Context, id uint) (*entity.Activity, error) {
	return s.actRepo.FindByID(ctx, id)
}

func (s *ActivityService) Create(ctx context.Context, a *entity.Activity) error {
	if err := s.actRepo.Create(ctx, a); err != nil {
		return err
	}
	// fan-out to followers' feeds
	followers, _ := s.followRepo.FindByFollowingID(ctx, a.ActorID)
	for _, f := range followers {
		feed := &entity.ActivityFeed{
			UserID:     f.FollowerID,
			ActorID:    a.ActorID,
			Type:       a.Type,
			Verb:       a.Verb,
			TargetType: a.TargetType,
			TargetID:   a.TargetID,
			Metadata:   a.Metadata,
			CreatedAt:  a.CreatedAt,
		}
		_ = s.feedRepo.Create(ctx, feed)
	}
	return nil
}

func (s *ActivityService) Delete(ctx context.Context, id uint) error {
	return s.actRepo.Delete(ctx, id)
}

func (s *ActivityService) Cleanup(ctx context.Context, days int) (int64, error) {
	return s.actRepo.DeleteOlderThan(ctx, days)
}

// Feed
func (s *ActivityService) GetFeed(ctx context.Context, userID uint, limit int) ([]*entity.ActivityFeed, error) {
	return s.feedRepo.FindByUserID(ctx, userID, limit)
}

func (s *ActivityService) GetUnreadCount(ctx context.Context, userID uint) (int64, error) {
	return s.feedRepo.CountUnread(ctx, userID)
}

func (s *ActivityService) MarkFeedRead(ctx context.Context, id uint) error {
	return s.feedRepo.MarkRead(ctx, id)
}

func (s *ActivityService) MarkAllFeedRead(ctx context.Context, userID uint) error {
	return s.feedRepo.MarkAllRead(ctx, userID)
}

func (s *ActivityService) CleanupFeed(ctx context.Context, days int) (int64, error) {
	return s.feedRepo.DeleteOlderThan(ctx, days)
}

// Follows
func (s *ActivityService) Follow(ctx context.Context, followerID, followingID uint) error {
	existing, _ := s.followRepo.FindAll(ctx, map[string]any{"follower_id": followerID, "following_id": followingID})
	if len(existing) > 0 {
		return nil // already following
	}
	f := &entity.Follow{FollowerID: followerID, FollowingID: followingID, CreatedAt: time.Now()}
	return s.followRepo.Create(ctx, f)
}

func (s *ActivityService) Unfollow(ctx context.Context, followerID, followingID uint) error {
	return s.followRepo.Delete(ctx, followerID, followingID)
}

func (s *ActivityService) IsFollowing(ctx context.Context, followerID, followingID uint) (bool, error) {
	return s.followRepo.IsFollowing(ctx, followerID, followingID)
}

func (s *ActivityService) GetFollowers(ctx context.Context, userID uint) ([]*entity.Follow, error) {
	return s.followRepo.FindByFollowingID(ctx, userID)
}

func (s *ActivityService) GetFollowing(ctx context.Context, userID uint) ([]*entity.Follow, error) {
	return s.followRepo.FindByFollowerID(ctx, userID)
}

func (s *ActivityService) GetFollowerCount(ctx context.Context, userID uint) (int64, error) {
	return s.followRepo.CountFollowers(ctx, userID)
}

func (s *ActivityService) GetFollowingCount(ctx context.Context, userID uint) (int64, error) {
	return s.followRepo.CountFollowing(ctx, userID)
}

// Engagement metrics
func (s *ActivityService) GetMetrics(ctx context.Context, userID uint) (*entity.EngagementMetric, error) {
	m, err := s.metricRepo.FindByUserID(ctx, userID)
	if err != nil {
		m = &entity.EngagementMetric{UserID: userID}
		_ = s.metricRepo.Upsert(ctx, m)
	}
	return m, nil
}

func (s *ActivityService) GetLeaderboard(ctx context.Context, limit int) ([]*entity.EngagementMetric, error) {
	return s.metricRepo.TopUsers(ctx, limit)
}

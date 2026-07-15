package service

import (
	"context"
	"errors"
	"go-modular/modules/analytics/domain/entity"
	"go-modular/modules/analytics/domain/repository"
	"time"
)

var (
	ErrEventNotFound = errors.New("event not found")
)

type AnalyticsService struct {
	eventRepo repository.AnalyticsEventRepository
	blogRepo  repository.BlogStatsRepository
	userRepo  repository.UserAnalyticsRepository
}

func NewAnalyticsService(eventRepo repository.AnalyticsEventRepository, blogRepo repository.BlogStatsRepository, userRepo repository.UserAnalyticsRepository) *AnalyticsService {
	return &AnalyticsService{eventRepo: eventRepo, blogRepo: blogRepo, userRepo: userRepo}
}

func (s *AnalyticsService) TrackEvent(ctx context.Context, e *entity.AnalyticsEvent) error {
	return s.eventRepo.Create(ctx, e)
}

func (s *AnalyticsService) ListEvents(ctx context.Context, filter map[string]any) ([]*entity.AnalyticsEvent, error) {
	return s.eventRepo.FindAll(ctx, filter)
}

func (s *AnalyticsService) CountByType(ctx context.Context, eventType string, since string) (int64, error) {
	return s.eventRepo.CountByEventType(ctx, eventType, since)
}

func (s *AnalyticsService) TrackPageView(ctx context.Context, path, ip, userAgent, referrer string, userID *uint) error {
	e := &entity.AnalyticsEvent{
		EventType: "page_view",
		UserID:    userID,
		Path:      path,
		IP:        ip,
		UserAgent: userAgent,
		Referrer:  referrer,
		CreatedAt: time.Now(),
	}
	return s.eventRepo.Create(ctx, e)
}

func (s *AnalyticsService) TrackBlogView(ctx context.Context, blogID uint, userID *uint) error {
	e := &entity.AnalyticsEvent{
		EventType:    "blog_view",
		UserID:       userID,
		ResourceType: "blog",
		ResourceID:   &blogID,
		Path:         "/blogs/",
		CreatedAt:    time.Now(),
	}
	if err := s.eventRepo.Create(ctx, e); err != nil {
		return err
	}
	// update blog stats
	return s.blogRepo.IncrementView(ctx, blogID)
}

func (s *AnalyticsService) GetBlogStats(ctx context.Context, blogID uint) (*entity.BlogStats, error) {
	stats, err := s.blogRepo.FindByBlogID(ctx, blogID)
	if err != nil {
		stats = &entity.BlogStats{BlogID: blogID, ViewCount: 0, CreatedAt: time.Now()}
		_ = s.blogRepo.Upsert(ctx, stats)
	}
	return stats, nil
}

func (s *AnalyticsService) GetTopBlogs(ctx context.Context, limit int) ([]*entity.BlogStats, error) {
	return s.blogRepo.TopBlogs(ctx, limit)
}

func (s *AnalyticsService) UpdateUserActivity(ctx context.Context, userID uint) error {
	return s.userRepo.UpdateLastActive(ctx, userID)
}

func (s *AnalyticsService) GetUserStats(ctx context.Context, userID uint) (*entity.UserAnalytics, error) {
	u, err := s.userRepo.FindByUserID(ctx, userID)
	if err != nil {
		u = &entity.UserAnalytics{UserID: userID, CreatedAt: time.Now()}
		_ = s.userRepo.Upsert(ctx, u)
	}
	return u, nil
}

func (s *AnalyticsService) Cleanup(ctx context.Context, days int) (int64, error) {
	return s.eventRepo.DeleteOlderThan(ctx, days)
}

// Dashboard aggregations
type DashboardStats struct {
	TotalViews  int64            `json:"total_views"`
	TotalUsers  int64            `json:"total_users"`
	TotalBlogs  int64            `json:"total_blogs"`
	TopBlogs    []*entity.BlogStats `json:"top_blogs"`
	RecentViews int64            `json:"recent_views"`
}

func (s *AnalyticsService) GetDashboardStats(ctx context.Context) (*DashboardStats, error) {
	views24h, _ := s.eventRepo.CountByEventType(ctx, "page_view", time.Now().Add(-24*time.Hour).Format(time.RFC3339))
	topBlogs, _ := s.blogRepo.TopBlogs(ctx, 5)
	return &DashboardStats{
		RecentViews: views24h,
		TopBlogs:    topBlogs,
	}, nil
}

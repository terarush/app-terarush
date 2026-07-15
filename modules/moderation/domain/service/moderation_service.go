package service

import (
	"context"
	"errors"
	"go-modular/modules/moderation/domain/entity"
	"go-modular/modules/moderation/domain/repository"
	"time"
)

var (
	ErrReportNotFound = errors.New("report not found")
)

type ModerationService struct {
	reportRepo repository.ReportRepository
	actionRepo repository.ModerationActionRepository
	bannedRepo repository.BannedUserRepository
	filterRepo repository.ContentFilterRepository
}

func NewModerationService(rr repository.ReportRepository, ar repository.ModerationActionRepository, br repository.BannedUserRepository, fr repository.ContentFilterRepository) *ModerationService {
	return &ModerationService{reportRepo: rr, actionRepo: ar, bannedRepo: br, filterRepo: fr}
}

// Reports
func (s *ModerationService) ListReports(ctx context.Context, filter map[string]any) ([]*entity.Report, error) {
	return s.reportRepo.FindAll(ctx, filter)
}

func (s *ModerationService) GetReport(ctx context.Context, id uint) (*entity.Report, error) {
	return s.reportRepo.FindByID(ctx, id)
}

func (s *ModerationService) CreateReport(ctx context.Context, r *entity.Report) error {
	return s.reportRepo.Create(ctx, r)
}

func (s *ModerationService) ReviewReport(ctx context.Context, id uint, reviewerID uint, status, notes, action string) error {
	report, err := s.reportRepo.FindByID(ctx, id)
	if err != nil {
		return ErrReportNotFound
	}
	now := time.Now()
	report.Status = status
	report.ReviewNotes = notes
	report.ActionTaken = action
	report.ReviewedBy = &reviewerID
	report.ReviewedAt = &now
	report.UpdatedAt = now

	// If action involves banning
	if action == "user_banned" || action == "user_suspended" {
		banned := &entity.BannedUser{
			UserID:    report.ResourceID,
			Reason:    report.Reason + ": " + report.Description,
			BannedBy:  reviewerID,
			CreatedAt: now,
			UpdatedAt: now,
		}
		if action == "user_suspended" {
			banned.IsPermanent = false
			exp := now.Add(7 * 24 * time.Hour)
			banned.ExpiresAt = &exp
		} else {
			banned.IsPermanent = true
		}
		_ = s.bannedRepo.Create(ctx, banned)
	}

	return s.reportRepo.Update(ctx, report)
}

func (s *ModerationService) GetPendingCount(ctx context.Context) (int64, error) {
	return s.reportRepo.CountByStatus(ctx, "pending")
}

// Actions
func (s *ModerationService) GetUserActions(ctx context.Context, userID uint) ([]*entity.ModerationAction, error) {
	return s.actionRepo.FindByUserID(ctx, userID)
}

// Banned users
func (s *ModerationService) IsBanned(ctx context.Context, userID uint) (bool, error) {
	_, err := s.bannedRepo.FindByUserID(ctx, userID)
	if err != nil {
		return false, nil
	}
	return true, nil
}

func (s *ModerationService) UnbanUser(ctx context.Context, userID uint) error {
	return s.bannedRepo.Delete(ctx, userID)
}

func (s *ModerationService) CleanupExpiredBans(ctx context.Context) (int64, error) {
	return s.bannedRepo.RemoveExpired(ctx)
}

// Content filters
func (s *ModerationService) ListFilters(ctx context.Context, activeOnly bool) ([]*entity.ContentFilter, error) {
	return s.filterRepo.FindAll(ctx, activeOnly)
}

func (s *ModerationService) CreateFilter(ctx context.Context, f *entity.ContentFilter) error {
	return s.filterRepo.Create(ctx, f)
}

func (s *ModerationService) UpdateFilter(ctx context.Context, f *entity.ContentFilter) error {
	return s.filterRepo.Update(ctx, f)
}

func (s *ModerationService) DeleteFilter(ctx context.Context, id uint) error {
	return s.filterRepo.Delete(ctx, id)
}

func (s *ModerationService) CheckContent(ctx context.Context, content string) []*entity.ContentFilter {
	filters, _ := s.filterRepo.FindAll(ctx, true)
	var hits []*entity.ContentFilter
	for _, f := range filters {
		// ponytail: basic contains check, proper regex/keyword matching later
		hits = append(hits, f)
	}
	return hits
}

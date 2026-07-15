package repository

import (
	"context"
	"errors"
	"go-modular/internal/pkg/database"
	"gorm.io/gorm"
	"go-modular/modules/moderation/domain/entity"
	"time"
)

var (
	ErrReportNotFound   = errors.New("report not found")
	ErrFilterNotFound   = errors.New("content filter not found")
	ErrBannedUserNotFound = errors.New("banned user not found")
)

type ReportRepositoryImpl struct{}

func (r ReportRepositoryImpl) FindAll(ctx context.Context, filter map[string]any) ([]*entity.Report, error) {
	var items []*entity.Report
	query := database.DB.WithContext(ctx)
	for k, v := range filter {
		query = query.Where(k, v)
	}
	result := query.Order("severity DESC, created_at ASC").Find(&items)
	return items, result.Error
}

func (r ReportRepositoryImpl) FindByID(ctx context.Context, id uint) (*entity.Report, error) {
	var item entity.Report
	err := database.DB.WithContext(ctx).First(&item, id).Error
	if err != nil {
		return nil, ErrReportNotFound
	}
	return &item, nil
}

func (r ReportRepositoryImpl) FindByStatus(ctx context.Context, status string, limit int) ([]*entity.Report, error) {
	var items []*entity.Report
	query := database.DB.WithContext(ctx).Where("status = ?", status).Order("severity DESC, created_at ASC")
	if limit > 0 {
		query = query.Limit(limit)
	}
	err := query.Find(&items).Error
	return items, err
}

func (r ReportRepositoryImpl) Create(ctx context.Context, report *entity.Report) error {
	return database.DB.WithContext(ctx).Create(report).Error
}

func (r ReportRepositoryImpl) Update(ctx context.Context, report *entity.Report) error {
	return database.DB.WithContext(ctx).Save(report).Error
}

func (r ReportRepositoryImpl) Delete(ctx context.Context, id uint) error {
	return database.DB.WithContext(ctx).Delete(&entity.Report{}, id).Error
}

func (r ReportRepositoryImpl) CountByStatus(ctx context.Context, status string) (int64, error) {
	var count int64
	err := database.DB.WithContext(ctx).Model(&entity.Report{}).Where("status = ?", status).Count(&count).Error
	return count, err
}

func NewReportRepository() ReportRepository {
	return ReportRepositoryImpl{}
}

type ModerationActionRepositoryImpl struct{}

func (r ModerationActionRepositoryImpl) FindAll(ctx context.Context, filter map[string]any) ([]*entity.ModerationAction, error) {
	var items []*entity.ModerationAction
	query := database.DB.WithContext(ctx)
	for k, v := range filter {
		query = query.Where(k, v)
	}
	err := query.Order("created_at DESC").Find(&items).Error
	return items, err
}

func (r ModerationActionRepositoryImpl) FindByUserID(ctx context.Context, userID uint) ([]*entity.ModerationAction, error) {
	var items []*entity.ModerationAction
	err := database.DB.WithContext(ctx).Where("user_id = ?", userID).Order("created_at DESC").Find(&items).Error
	return items, err
}

func (r ModerationActionRepositoryImpl) Create(ctx context.Context, a *entity.ModerationAction) error {
	return database.DB.WithContext(ctx).Create(a).Error
}

func NewModerationActionRepository() ModerationActionRepository {
	return ModerationActionRepositoryImpl{}
}

type BannedUserRepositoryImpl struct{}

func (r BannedUserRepositoryImpl) FindByUserID(ctx context.Context, userID uint) (*entity.BannedUser, error) {
	var b entity.BannedUser
	err := database.DB.WithContext(ctx).Where("user_id = ?", userID).First(&b).Error
	if err != nil {
		return nil, ErrBannedUserNotFound
	}
	return &b, nil
}

func (r BannedUserRepositoryImpl) FindAllActive(ctx context.Context) ([]*entity.BannedUser, error) {
	var items []*entity.BannedUser
	err := database.DB.WithContext(ctx).
		Where("is_permanent = ? OR expires_at > ?", true, time.Now()).
		Find(&items).Error
	return items, err
}

func (r BannedUserRepositoryImpl) Create(ctx context.Context, b *entity.BannedUser) error {
	return database.DB.WithContext(ctx).Create(b).Error
}

func (r BannedUserRepositoryImpl) Delete(ctx context.Context, userID uint) error {
	return database.DB.WithContext(ctx).Where("user_id = ?", userID).Delete(&entity.BannedUser{}).Error
}

func (r BannedUserRepositoryImpl) RemoveExpired(ctx context.Context) (int64, error) {
	result := database.DB.WithContext(ctx).
		Where("is_permanent = ? AND expires_at < ?", false, time.Now()).
		Delete(&entity.BannedUser{})
	return result.RowsAffected, result.Error
}

func NewBannedUserRepository() BannedUserRepository {
	return BannedUserRepositoryImpl{}
}

type ContentFilterRepositoryImpl struct{}

func (r ContentFilterRepositoryImpl) FindAll(ctx context.Context, activeOnly bool) ([]*entity.ContentFilter, error) {
	var items []*entity.ContentFilter
	query := database.DB.WithContext(ctx)
	if activeOnly {
		query = query.Where("is_active = ?", true)
	}
	err := query.Order("type, severity DESC").Find(&items).Error
	return items, err
}

func (r ContentFilterRepositoryImpl) FindByType(ctx context.Context, filterType string) ([]*entity.ContentFilter, error) {
	var items []*entity.ContentFilter
	err := database.DB.WithContext(ctx).Where("type = ? AND is_active = ?", filterType, true).Find(&items).Error
	return items, err
}

func (r ContentFilterRepositoryImpl) FindByID(ctx context.Context, id uint) (*entity.ContentFilter, error) {
	var f entity.ContentFilter
	err := database.DB.WithContext(ctx).First(&f, id).Error
	if err != nil {
		return nil, ErrFilterNotFound
	}
	return &f, nil
}

func (r ContentFilterRepositoryImpl) Create(ctx context.Context, f *entity.ContentFilter) error {
	return database.DB.WithContext(ctx).Create(f).Error
}

func (r ContentFilterRepositoryImpl) Update(ctx context.Context, f *entity.ContentFilter) error {
	return database.DB.WithContext(ctx).Save(f).Error
}

func (r ContentFilterRepositoryImpl) Delete(ctx context.Context, id uint) error {
	return database.DB.WithContext(ctx).Delete(&entity.ContentFilter{}, id).Error
}

func (r ContentFilterRepositoryImpl) IncrementHit(ctx context.Context, id uint) error {
	return database.DB.WithContext(ctx).Model(&entity.ContentFilter{}).
		Where("id = ?", id).
		UpdateColumn("hit_count", gorm.Expr("hit_count + 1")).Error
}

func NewContentFilterRepository() ContentFilterRepository {
	return ContentFilterRepositoryImpl{}
}

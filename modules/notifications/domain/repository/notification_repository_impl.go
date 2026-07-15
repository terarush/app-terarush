package repository

import (
	"context"
	"errors"
	"go-modular/internal/pkg/database"
	"go-modular/modules/notifications/domain/entity"
	"time"
)

var (
	ErrNotificationNotFound = errors.New("notification not found")
)

type NotificationRepositoryImpl struct{}

func (r NotificationRepositoryImpl) FindAll(ctx context.Context, filter map[string]any) ([]*entity.Notification, error) {
	var items []*entity.Notification
	query := database.DB.WithContext(ctx)
	for k, v := range filter {
		query = query.Where(k, v)
	}
	result := query.Order("created_at DESC").Find(&items)
	return items, result.Error
}

func (r NotificationRepositoryImpl) FindByID(ctx context.Context, id uint) (*entity.Notification, error) {
	var n entity.Notification
	result := database.DB.WithContext(ctx).First(&n, id)
	if result.Error != nil {
		return nil, ErrNotificationNotFound
	}
	return &n, nil
}

func (r NotificationRepositoryImpl) FindByUserID(ctx context.Context, userID uint, status string, limit int) ([]*entity.Notification, error) {
	var items []*entity.Notification
	query := database.DB.WithContext(ctx).Where("user_id = ?", userID)
	if status == "unread" {
		query = query.Where("is_read = ?", false)
	} else if status == "archived" {
		query = query.Where("is_archived = ?", true)
	}
	if limit > 0 {
		query = query.Limit(limit)
	}
	result := query.Order("created_at DESC").Find(&items)
	return items, result.Error
}

func (r NotificationRepositoryImpl) CountUnread(ctx context.Context, userID uint) (int64, error) {
	var count int64
	err := database.DB.WithContext(ctx).Model(&entity.Notification{}).
		Where("user_id = ? AND is_read = ? AND is_archived = ?", userID, false, false).
		Count(&count).Error
	return count, err
}

func (r NotificationRepositoryImpl) MarkRead(ctx context.Context, id uint) error {
	now := time.Now()
	return database.DB.WithContext(ctx).Model(&entity.Notification{}).
		Where("id = ?", id).
		Updates(map[string]any{"is_read": true, "read_at": now}).Error
}

func (r NotificationRepositoryImpl) MarkAllRead(ctx context.Context, userID uint) error {
	now := time.Now()
	return database.DB.WithContext(ctx).Model(&entity.Notification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Updates(map[string]any{"is_read": true, "read_at": now}).Error
}

func (r NotificationRepositoryImpl) Create(ctx context.Context, n *entity.Notification) error {
	return database.DB.WithContext(ctx).Create(n).Error
}

func (r NotificationRepositoryImpl) Update(ctx context.Context, n *entity.Notification) error {
	return database.DB.WithContext(ctx).Save(n).Error
}

func (r NotificationRepositoryImpl) Delete(ctx context.Context, id uint) error {
	return database.DB.WithContext(ctx).Delete(&entity.Notification{}, id).Error
}

func (r NotificationRepositoryImpl) DeleteOlderThan(ctx context.Context, days int) (int64, error) {
	cutoff := time.Now().AddDate(0, 0, -days)
	result := database.DB.WithContext(ctx).Where("created_at < ?", cutoff).Delete(&entity.Notification{})
	return result.RowsAffected, result.Error
}

func NewNotificationRepository() NotificationRepository {
	return NotificationRepositoryImpl{}
}

type NotificationTemplateRepositoryImpl struct{}

func (r NotificationTemplateRepositoryImpl) FindAll(ctx context.Context) ([]*entity.NotificationTemplate, error) {
	var items []*entity.NotificationTemplate
	err := database.DB.WithContext(ctx).Find(&items).Error
	return items, err
}

func (r NotificationTemplateRepositoryImpl) FindByID(ctx context.Context, id uint) (*entity.NotificationTemplate, error) {
	var t entity.NotificationTemplate
	err := database.DB.WithContext(ctx).First(&t, id).Error
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r NotificationTemplateRepositoryImpl) FindByName(ctx context.Context, name string) (*entity.NotificationTemplate, error) {
	var t entity.NotificationTemplate
	err := database.DB.WithContext(ctx).Where("name = ?", name).First(&t).Error
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r NotificationTemplateRepositoryImpl) Create(ctx context.Context, t *entity.NotificationTemplate) error {
	return database.DB.WithContext(ctx).Create(t).Error
}

func (r NotificationTemplateRepositoryImpl) Update(ctx context.Context, t *entity.NotificationTemplate) error {
	return database.DB.WithContext(ctx).Save(t).Error
}

func (r NotificationTemplateRepositoryImpl) Delete(ctx context.Context, id uint) error {
	return database.DB.WithContext(ctx).Delete(&entity.NotificationTemplate{}, id).Error
}

func NewNotificationTemplateRepository() NotificationTemplateRepository {
	return NotificationTemplateRepositoryImpl{}
}

type NotificationPreferenceRepositoryImpl struct{}

func (r NotificationPreferenceRepositoryImpl) FindByUserID(ctx context.Context, userID uint) (*entity.NotificationPreference, error) {
	var p entity.NotificationPreference
	err := database.DB.WithContext(ctx).Where("user_id = ?", userID).First(&p).Error
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r NotificationPreferenceRepositoryImpl) Upsert(ctx context.Context, p *entity.NotificationPreference) error {
	return database.DB.WithContext(ctx).Save(p).Error
}

func NewNotificationPreferenceRepository() NotificationPreferenceRepository {
	return NotificationPreferenceRepositoryImpl{}
}

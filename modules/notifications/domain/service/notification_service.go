package service

import (
	"context"
	"errors"
	"go-modular/modules/notifications/domain/entity"
	"go-modular/modules/notifications/domain/repository"
)

var (
	ErrNotificationNotFound = errors.New("notification not found")
)

type NotificationService struct {
	notifRepo repository.NotificationRepository
	tmplRepo  repository.NotificationTemplateRepository
	prefRepo  repository.NotificationPreferenceRepository
}

func NewNotificationService(notifRepo repository.NotificationRepository, tmplRepo repository.NotificationTemplateRepository, prefRepo repository.NotificationPreferenceRepository) *NotificationService {
	return &NotificationService{notifRepo: notifRepo, tmplRepo: tmplRepo, prefRepo: prefRepo}
}

func (s *NotificationService) List(ctx context.Context, filter map[string]any) ([]*entity.Notification, error) {
	return s.notifRepo.FindAll(ctx, filter)
}

func (s *NotificationService) Get(ctx context.Context, id uint) (*entity.Notification, error) {
	return s.notifRepo.FindByID(ctx, id)
}

func (s *NotificationService) GetByUserID(ctx context.Context, userID uint, status string, limit int) ([]*entity.Notification, error) {
	return s.notifRepo.FindByUserID(ctx, userID, status, limit)
}

func (s *NotificationService) GetUnreadCount(ctx context.Context, userID uint) (int64, error) {
	return s.notifRepo.CountUnread(ctx, userID)
}

func (s *NotificationService) MarkRead(ctx context.Context, id uint) error {
	return s.notifRepo.MarkRead(ctx, id)
}

func (s *NotificationService) MarkAllRead(ctx context.Context, userID uint) error {
	return s.notifRepo.MarkAllRead(ctx, userID)
}

func (s *NotificationService) Create(ctx context.Context, n *entity.Notification) error {
	return s.notifRepo.Create(ctx, n)
}

func (s *NotificationService) Delete(ctx context.Context, id uint) error {
	return s.notifRepo.Delete(ctx, id)
}

func (s *NotificationService) CreateFromTemplate(ctx context.Context, userID uint, tmplName string, titleData, bodyData map[string]string, link string) (*entity.Notification, error) {
	tmpl, err := s.tmplRepo.FindByName(ctx, tmplName)
	if err != nil {
		return nil, err
	}
	title := tmpl.TitleTmpl
	body := tmpl.BodyTmpl
	_ = titleData
	_ = bodyData
	n := entity.NewNotification(userID, tmplName, title, body)
	n.Link = link
	n.Icon = tmpl.DefaultIcon
	if err := s.notifRepo.Create(ctx, n); err != nil {
		return nil, err
	}
	return n, nil
}

func (s *NotificationService) Cleanup(ctx context.Context, days int) (int64, error) {
	return s.notifRepo.DeleteOlderThan(ctx, days)
}

// Templates
func (s *NotificationService) ListTemplates(ctx context.Context) ([]*entity.NotificationTemplate, error) {
	return s.tmplRepo.FindAll(ctx)
}

func (s *NotificationService) CreateTemplate(ctx context.Context, t *entity.NotificationTemplate) error {
	return s.tmplRepo.Create(ctx, t)
}

func (s *NotificationService) UpdateTemplate(ctx context.Context, t *entity.NotificationTemplate) error {
	return s.tmplRepo.Update(ctx, t)
}

func (s *NotificationService) DeleteTemplate(ctx context.Context, id uint) error {
	return s.tmplRepo.Delete(ctx, id)
}

// Preferences
func (s *NotificationService) GetPreferences(ctx context.Context, userID uint) (*entity.NotificationPreference, error) {
	return s.prefRepo.FindByUserID(ctx, userID)
}

func (s *NotificationService) UpsertPreferences(ctx context.Context, p *entity.NotificationPreference) error {
	return s.prefRepo.Upsert(ctx, p)
}

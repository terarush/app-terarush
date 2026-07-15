package repository

import (
	"context"
	"errors"
	"go-modular/internal/pkg/database"
	"go-modular/modules/subscriptions/domain/entity"
)

var (
	ErrSubscriptionNotFound = errors.New("subscription not found")
	ErrNewsletterNotFound   = errors.New("newsletter not found")
)

type SubscriptionRepositoryImpl struct{}

func (r SubscriptionRepositoryImpl) FindAll(ctx context.Context, filter map[string]any) ([]*entity.Subscription, error) {
	var items []*entity.Subscription
	query := database.DB.WithContext(ctx)
	for k, v := range filter {
		query = query.Where(k, v)
	}
	return items, query.Order("created_at DESC").Find(&items).Error
}

func (r SubscriptionRepositoryImpl) FindByID(ctx context.Context, id uint) (*entity.Subscription, error) {
	var s entity.Subscription
	err := database.DB.WithContext(ctx).First(&s, id).Error
	if err != nil {
		return nil, ErrSubscriptionNotFound
	}
	return &s, nil
}

func (r SubscriptionRepositoryImpl) FindByUserAndResource(ctx context.Context, userID uint, resourceType string, resourceID uint) (*entity.Subscription, error) {
	var s entity.Subscription
	err := database.DB.WithContext(ctx).Where("user_id = ? AND resource_type = ? AND resource_id = ?", userID, resourceType, resourceID).First(&s).Error
	if err != nil {
		return nil, ErrSubscriptionNotFound
	}
	return &s, nil
}

func (r SubscriptionRepositoryImpl) FindByResource(ctx context.Context, resourceType string, resourceID uint, status string) ([]*entity.Subscription, error) {
	var items []*entity.Subscription
	query := database.DB.WithContext(ctx).Where("resource_type = ? AND resource_id = ?", resourceType, resourceID)
	if status != "" {
		query = query.Where("status = ?", status)
	}
	return items, query.Find(&items).Error
}

func (r SubscriptionRepositoryImpl) FindByUserID(ctx context.Context, userID uint) ([]*entity.Subscription, error) {
	var items []*entity.Subscription
	err := database.DB.WithContext(ctx).Where("user_id = ?", userID).Order("created_at DESC").Find(&items).Error
	return items, err
}

func (r SubscriptionRepositoryImpl) Create(ctx context.Context, s *entity.Subscription) error {
	return database.DB.WithContext(ctx).Create(s).Error
}

func (r SubscriptionRepositoryImpl) Update(ctx context.Context, s *entity.Subscription) error {
	return database.DB.WithContext(ctx).Save(s).Error
}

func (r SubscriptionRepositoryImpl) Delete(ctx context.Context, id uint) error {
	return database.DB.WithContext(ctx).Delete(&entity.Subscription{}, id).Error
}

func (r SubscriptionRepositoryImpl) CountByResource(ctx context.Context, resourceType string, resourceID uint) (int64, error) {
	var count int64
	err := database.DB.WithContext(ctx).Model(&entity.Subscription{}).
		Where("resource_type = ? AND resource_id = ? AND status = ?", resourceType, resourceID, "active").
		Count(&count).Error
	return count, err
}

func NewSubscriptionRepository() SubscriptionRepository {
	return SubscriptionRepositoryImpl{}
}

type NewsletterRepositoryImpl struct{}

func (r NewsletterRepositoryImpl) FindAll(ctx context.Context, filter map[string]any) ([]*entity.Newsletter, error) {
	var items []*entity.Newsletter
	query := database.DB.WithContext(ctx)
	for k, v := range filter {
		query = query.Where(k, v)
	}
	return items, query.Order("created_at DESC").Find(&items).Error
}

func (r NewsletterRepositoryImpl) FindByID(ctx context.Context, id uint) (*entity.Newsletter, error) {
	var n entity.Newsletter
	err := database.DB.WithContext(ctx).First(&n, id).Error
	if err != nil {
		return nil, ErrNewsletterNotFound
	}
	return &n, nil
}

func (r NewsletterRepositoryImpl) FindByStatus(ctx context.Context, status string) ([]*entity.Newsletter, error) {
	var items []*entity.Newsletter
	err := database.DB.WithContext(ctx).Where("status = ?", status).Find(&items).Error
	return items, err
}

func (r NewsletterRepositoryImpl) Create(ctx context.Context, n *entity.Newsletter) error {
	return database.DB.WithContext(ctx).Create(n).Error
}

func (r NewsletterRepositoryImpl) Update(ctx context.Context, n *entity.Newsletter) error {
	return database.DB.WithContext(ctx).Save(n).Error
}

func (r NewsletterRepositoryImpl) Delete(ctx context.Context, id uint) error {
	return database.DB.WithContext(ctx).Delete(&entity.Newsletter{}, id).Error
}

func NewNewsletterRepository() NewsletterRepository {
	return NewsletterRepositoryImpl{}
}

type NewsletterSubscriberRepositoryImpl struct{}

func (r NewsletterSubscriberRepositoryImpl) FindAll(ctx context.Context, status string, limit int) ([]*entity.NewsletterSubscriber, error) {
	var items []*entity.NewsletterSubscriber
	query := database.DB.WithContext(ctx)
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if limit > 0 {
		query = query.Limit(limit)
	}
	return items, query.Find(&items).Error
}

func (r NewsletterSubscriberRepositoryImpl) FindByEmail(ctx context.Context, email string) (*entity.NewsletterSubscriber, error) {
	var s entity.NewsletterSubscriber
	err := database.DB.WithContext(ctx).Where("email = ?", email).First(&s).Error
	return &s, err
}

func (r NewsletterSubscriberRepositoryImpl) FindByToken(ctx context.Context, token string) (*entity.NewsletterSubscriber, error) {
	var s entity.NewsletterSubscriber
	err := database.DB.WithContext(ctx).Where("token = ?", token).First(&s).Error
	return &s, err
}

func (r NewsletterSubscriberRepositoryImpl) Create(ctx context.Context, s *entity.NewsletterSubscriber) error {
	return database.DB.WithContext(ctx).Create(s).Error
}

func (r NewsletterSubscriberRepositoryImpl) Update(ctx context.Context, s *entity.NewsletterSubscriber) error {
	return database.DB.WithContext(ctx).Save(s).Error
}

func (r NewsletterSubscriberRepositoryImpl) CountByStatus(ctx context.Context, status string) (int64, error) {
	var count int64
	err := database.DB.WithContext(ctx).Model(&entity.NewsletterSubscriber{}).
		Where("status = ?", status).Count(&count).Error
	return count, err
}

func NewNewsletterSubscriberRepository() NewsletterSubscriberRepository {
	return NewsletterSubscriberRepositoryImpl{}
}

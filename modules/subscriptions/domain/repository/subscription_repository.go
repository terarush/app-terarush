package repository

import (
	"context"
	"go-modular/modules/subscriptions/domain/entity"
)

type SubscriptionRepository interface {
	FindAll(ctx context.Context, filter map[string]any) ([]*entity.Subscription, error)
	FindByID(ctx context.Context, id uint) (*entity.Subscription, error)
	FindByUserAndResource(ctx context.Context, userID uint, resourceType string, resourceID uint) (*entity.Subscription, error)
	FindByResource(ctx context.Context, resourceType string, resourceID uint, status string) ([]*entity.Subscription, error)
	FindByUserID(ctx context.Context, userID uint) ([]*entity.Subscription, error)
	Create(ctx context.Context, s *entity.Subscription) error
	Update(ctx context.Context, s *entity.Subscription) error
	Delete(ctx context.Context, id uint) error
	CountByResource(ctx context.Context, resourceType string, resourceID uint) (int64, error)
}

type NewsletterRepository interface {
	FindAll(ctx context.Context, filter map[string]any) ([]*entity.Newsletter, error)
	FindByID(ctx context.Context, id uint) (*entity.Newsletter, error)
	FindByStatus(ctx context.Context, status string) ([]*entity.Newsletter, error)
	Create(ctx context.Context, n *entity.Newsletter) error
	Update(ctx context.Context, n *entity.Newsletter) error
	Delete(ctx context.Context, id uint) error
}

type NewsletterSubscriberRepository interface {
	FindAll(ctx context.Context, status string, limit int) ([]*entity.NewsletterSubscriber, error)
	FindByEmail(ctx context.Context, email string) (*entity.NewsletterSubscriber, error)
	FindByToken(ctx context.Context, token string) (*entity.NewsletterSubscriber, error)
	Create(ctx context.Context, s *entity.NewsletterSubscriber) error
	Update(ctx context.Context, s *entity.NewsletterSubscriber) error
	CountByStatus(ctx context.Context, status string) (int64, error)
}

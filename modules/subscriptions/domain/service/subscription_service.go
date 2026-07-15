package service

import (
	"context"
	"errors"
	"go-modular/modules/subscriptions/domain/entity"
	"go-modular/modules/subscriptions/domain/repository"
	"time"
)

var (
	ErrSubscriptionNotFound = errors.New("subscription not found")
)

type SubscriptionService struct {
	subRepo  repository.SubscriptionRepository
	newsRepo repository.NewsletterRepository
	subsRepo repository.NewsletterSubscriberRepository
}

func NewSubscriptionService(subRepo repository.SubscriptionRepository, newsRepo repository.NewsletterRepository, subsRepo repository.NewsletterSubscriberRepository) *SubscriptionService {
	return &SubscriptionService{subRepo: subRepo, newsRepo: newsRepo, subsRepo: subsRepo}
}

// Subscriptions
func (s *SubscriptionService) List(ctx context.Context, filter map[string]any) ([]*entity.Subscription, error) {
	return s.subRepo.FindAll(ctx, filter)
}

func (s *SubscriptionService) Get(ctx context.Context, id uint) (*entity.Subscription, error) {
	return s.subRepo.FindByID(ctx, id)
}

func (s *SubscriptionService) GetUserSubscriptions(ctx context.Context, userID uint) ([]*entity.Subscription, error) {
	return s.subRepo.FindByUserID(ctx, userID)
}

func (s *SubscriptionService) Subscribe(ctx context.Context, userID uint, resourceType string, resourceID uint) (*entity.Subscription, error) {
	existing, err := s.subRepo.FindByUserAndResource(ctx, userID, resourceType, resourceID)
	if err == nil && existing != nil {
		if existing.Status == "cancelled" {
			existing.Status = "active"
			existing.UpdatedAt = time.Now()
			_ = s.subRepo.Update(ctx, existing)
		}
		return existing, nil
	}
	sub := entity.NewSubscription(userID, resourceType, resourceID)
	return sub, s.subRepo.Create(ctx, sub)
}

func (s *SubscriptionService) Unsubscribe(ctx context.Context, id uint) error {
	sub, err := s.subRepo.FindByID(ctx, id)
	if err != nil {
		return ErrSubscriptionNotFound
	}
	sub.Status = "cancelled"
	sub.UpdatedAt = time.Now()
	return s.subRepo.Update(ctx, sub)
}

func (s *SubscriptionService) Delete(ctx context.Context, id uint) error {
	return s.subRepo.Delete(ctx, id)
}

func (s *SubscriptionService) GetSubscriberCount(ctx context.Context, resourceType string, resourceID uint) (int64, error) {
	return s.subRepo.CountByResource(ctx, resourceType, resourceID)
}

// Newsletters
func (s *SubscriptionService) ListNewsletters(ctx context.Context, filter map[string]any) ([]*entity.Newsletter, error) {
	return s.newsRepo.FindAll(ctx, filter)
}

func (s *SubscriptionService) GetNewsletter(ctx context.Context, id uint) (*entity.Newsletter, error) {
	return s.newsRepo.FindByID(ctx, id)
}

func (s *SubscriptionService) CreateNewsletter(ctx context.Context, n *entity.Newsletter) error {
	return s.newsRepo.Create(ctx, n)
}

func (s *SubscriptionService) SendNewsletter(ctx context.Context, id uint) error {
	n, err := s.newsRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	n.Status = "sending"
	n.SentAt = &[]time.Time{time.Now()}[0]
	subs, _ := s.subsRepo.FindAll(ctx, "active", 0)
	n.SentCount = len(subs)
	n.Status = "sent"
	return s.newsRepo.Update(ctx, n)
}

func (s *SubscriptionService) DeleteNewsletter(ctx context.Context, id uint) error {
	return s.newsRepo.Delete(ctx, id)
}

// Subscribers
func (s *SubscriptionService) ListSubscribers(ctx context.Context, status string, limit int) ([]*entity.NewsletterSubscriber, error) {
	return s.subsRepo.FindAll(ctx, status, limit)
}

func (s *SubscriptionService) AddSubscriber(ctx context.Context, email, name string) (*entity.NewsletterSubscriber, error) {
	sub := &entity.NewsletterSubscriber{
		Email:     email,
		Name:      name,
		Status:    "active",
		Token:     "", // generate UUID
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	return sub, s.subsRepo.Create(ctx, sub)
}

func (s *SubscriptionService) RemoveSubscriber(ctx context.Context, email string) error {
	sub, err := s.subsRepo.FindByEmail(ctx, email)
	if err != nil {
		return err
	}
	sub.Status = "unsubscribed"
	return s.subsRepo.Update(ctx, sub)
}

func (s *SubscriptionService) GetSubscriberCountByStatus(ctx context.Context, status string) (int64, error) {
	return s.subsRepo.CountByStatus(ctx, status)
}

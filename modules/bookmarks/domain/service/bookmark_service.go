package service

import (
	"context"
	"errors"
	"go-modular/modules/bookmarks/domain/entity"
	"go-modular/modules/bookmarks/domain/repository"
	"time"
)

var (
	ErrBookmarkNotFound   = errors.New("bookmark not found")
	ErrCollectionNotFound = errors.New("collection not found")
)

type BookmarkService struct {
	bmRepo   repository.BookmarkRepository
	collRepo repository.CollectionRepository
	goalRepo repository.ReadingGoalRepository
}

func NewBookmarkService(bmRepo repository.BookmarkRepository, collRepo repository.CollectionRepository, goalRepo repository.ReadingGoalRepository) *BookmarkService {
	return &BookmarkService{bmRepo: bmRepo, collRepo: collRepo, goalRepo: goalRepo}
}

// Bookmarks
func (s *BookmarkService) List(ctx context.Context, filter map[string]any) ([]*entity.Bookmark, error) {
	return s.bmRepo.FindAll(ctx, filter)
}

func (s *BookmarkService) GetUserBookmarks(ctx context.Context, userID uint, status string) ([]*entity.Bookmark, error) {
	return s.bmRepo.FindByUserID(ctx, userID, status)
}

func (s *BookmarkService) Get(ctx context.Context, id uint) (*entity.Bookmark, error) {
	return s.bmRepo.FindByID(ctx, id)
}

func (s *BookmarkService) Add(ctx context.Context, userID uint, resourceType string, resourceID uint, collectionID *uint, note, tags string) (*entity.Bookmark, error) {
	bm := entity.NewBookmark(userID, resourceType, resourceID)
	bm.CollectionID = collectionID
	bm.Note = note
	bm.Tags = tags
	return bm, s.bmRepo.Create(ctx, bm)
}

func (s *BookmarkService) UpdateProgress(ctx context.Context, id uint, progress float64, status string) error {
	bm, err := s.bmRepo.FindByID(ctx, id)
	if err != nil {
		return ErrBookmarkNotFound
	}
	bm.Progress = progress
	bm.UpdatedAt = time.Now()
	if status != "" {
		bm.Status = status
		// increment reading goal if completed
		if status == "completed" {
			_ = s.goalRepo.Increment(ctx, bm.UserID, time.Now().Year())
		}
	}
	return s.bmRepo.Update(ctx, bm)
}

func (s *BookmarkService) Remove(ctx context.Context, id uint) error {
	return s.bmRepo.Delete(ctx, id)
}

// Collections
func (s *BookmarkService) ListCollections(ctx context.Context, userID uint) ([]*entity.Collection, error) {
	return s.collRepo.FindAll(ctx, userID)
}

func (s *BookmarkService) CreateCollection(ctx context.Context, c *entity.Collection) error {
	return s.collRepo.Create(ctx, c)
}

func (s *BookmarkService) UpdateCollection(ctx context.Context, c *entity.Collection) error {
	return s.collRepo.Update(ctx, c)
}

func (s *BookmarkService) DeleteCollection(ctx context.Context, id uint) error {
	return s.collRepo.Delete(ctx, id)
}

// Reading Goals
func (s *BookmarkService) GetGoal(ctx context.Context, userID uint, year int) (*entity.ReadingGoal, error) {
	g, err := s.goalRepo.FindByUserAndYear(ctx, userID, year)
	if err != nil {
		g = entity.NewReadingGoal(userID, year, 0)
		_ = s.goalRepo.Upsert(ctx, g)
	}
	return g, nil
}

func (s *BookmarkService) SetGoal(ctx context.Context, userID uint, year, goalCount int) (*entity.ReadingGoal, error) {
	g, err := s.goalRepo.FindByUserAndYear(ctx, userID, year)
	if err != nil {
		g = entity.NewReadingGoal(userID, year, goalCount)
		return g, s.goalRepo.Upsert(ctx, g)
	}
	g.Goal = goalCount
	g.UpdatedAt = time.Now()
	return g, s.goalRepo.Upsert(ctx, g)
}

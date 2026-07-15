package repository

import (
	"context"
	"errors"
	"go-modular/internal/pkg/database"
	"gorm.io/gorm"
	"go-modular/modules/bookmarks/domain/entity"
)

var (
	ErrBookmarkNotFound   = errors.New("bookmark not found")
	ErrCollectionNotFound = errors.New("collection not found")
	ErrGoalNotFound       = errors.New("reading goal not found")
)

type BookmarkRepositoryImpl struct{}

func (r BookmarkRepositoryImpl) FindAll(ctx context.Context, filter map[string]any) ([]*entity.Bookmark, error) {
	var items []*entity.Bookmark
	query := database.DB.WithContext(ctx)
	for k, v := range filter {
		query = query.Where(k, v)
	}
	return items, query.Order("created_at DESC").Find(&items).Error
}

func (r BookmarkRepositoryImpl) FindByID(ctx context.Context, id uint) (*entity.Bookmark, error) {
	var b entity.Bookmark
	err := database.DB.WithContext(ctx).First(&b, id).Error
	if err != nil {
		return nil, ErrBookmarkNotFound
	}
	return &b, nil
}

func (r BookmarkRepositoryImpl) FindByUserID(ctx context.Context, userID uint, status string) ([]*entity.Bookmark, error) {
	var items []*entity.Bookmark
	query := database.DB.WithContext(ctx).Where("user_id = ?", userID)
	if status != "" {
		query = query.Where("status = ?", status)
	}
	return items, query.Order("created_at DESC").Find(&items).Error
}

func (r BookmarkRepositoryImpl) FindByUserAndResource(ctx context.Context, userID uint, resourceType string, resourceID uint) (*entity.Bookmark, error) {
	var b entity.Bookmark
	err := database.DB.WithContext(ctx).Where("user_id = ? AND resource_type = ? AND resource_id = ?", userID, resourceType, resourceID).First(&b).Error
	if err != nil {
		return nil, ErrBookmarkNotFound
	}
	return &b, nil
}

func (r BookmarkRepositoryImpl) FindByCollectionID(ctx context.Context, collectionID uint) ([]*entity.Bookmark, error) {
	var items []*entity.Bookmark
	err := database.DB.WithContext(ctx).Where("collection_id = ?", collectionID).Order("created_at DESC").Find(&items).Error
	return items, err
}

func (r BookmarkRepositoryImpl) Create(ctx context.Context, b *entity.Bookmark) error {
	return database.DB.WithContext(ctx).Create(b).Error
}

func (r BookmarkRepositoryImpl) Update(ctx context.Context, b *entity.Bookmark) error {
	return database.DB.WithContext(ctx).Save(b).Error
}

func (r BookmarkRepositoryImpl) Delete(ctx context.Context, id uint) error {
	return database.DB.WithContext(ctx).Delete(&entity.Bookmark{}, id).Error
}

func (r BookmarkRepositoryImpl) CountByUserID(ctx context.Context, userID uint) (int64, error) {
	var count int64
	err := database.DB.WithContext(ctx).Model(&entity.Bookmark{}).Where("user_id = ?", userID).Count(&count).Error
	return count, err
}

func NewBookmarkRepository() BookmarkRepository {
	return BookmarkRepositoryImpl{}
}

type CollectionRepositoryImpl struct{}

func (r CollectionRepositoryImpl) FindAll(ctx context.Context, userID uint) ([]*entity.Collection, error) {
	var items []*entity.Collection
	err := database.DB.WithContext(ctx).Where("user_id = ?", userID).Order("sort_order ASC, created_at ASC").Find(&items).Error
	return items, err
}

func (r CollectionRepositoryImpl) FindByID(ctx context.Context, id uint) (*entity.Collection, error) {
	var c entity.Collection
	err := database.DB.WithContext(ctx).First(&c, id).Error
	if err != nil {
		return nil, ErrCollectionNotFound
	}
	return &c, nil
}

func (r CollectionRepositoryImpl) Create(ctx context.Context, c *entity.Collection) error {
	return database.DB.WithContext(ctx).Create(c).Error
}

func (r CollectionRepositoryImpl) Update(ctx context.Context, c *entity.Collection) error {
	return database.DB.WithContext(ctx).Save(c).Error
}

func (r CollectionRepositoryImpl) Delete(ctx context.Context, id uint) error {
	return database.DB.WithContext(ctx).Delete(&entity.Collection{}, id).Error
}

func NewCollectionRepository() CollectionRepository {
	return CollectionRepositoryImpl{}
}

type ReadingGoalRepositoryImpl struct{}

func (r ReadingGoalRepositoryImpl) FindByUserAndYear(ctx context.Context, userID uint, year int) (*entity.ReadingGoal, error) {
	var g entity.ReadingGoal
	err := database.DB.WithContext(ctx).Where("user_id = ? AND year = ?", userID, year).First(&g).Error
	if err != nil {
		return nil, ErrGoalNotFound
	}
	return &g, nil
}

func (r ReadingGoalRepositoryImpl) Upsert(ctx context.Context, g *entity.ReadingGoal) error {
	return database.DB.WithContext(ctx).Save(g).Error
}

func (r ReadingGoalRepositoryImpl) Increment(ctx context.Context, userID uint, year int) error {
	return database.DB.WithContext(ctx).Model(&entity.ReadingGoal{}).
		Where("user_id = ? AND year = ?", userID, year).
		UpdateColumn("progress", gorm.Expr("progress + 1")).Error
}

func NewReadingGoalRepository() ReadingGoalRepository {
	return ReadingGoalRepositoryImpl{}
}

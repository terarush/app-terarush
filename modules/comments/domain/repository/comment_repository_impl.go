package repository

import (
	"context"
	"go-modular/internal/pkg/database"
	"go-modular/modules/comments/domain/entity"
)

var (
	ERR_RECORD_NOT_FOUND = "record not found"
)

type CommentRepositoryImpl struct{}

func (r *CommentRepositoryImpl) Create(ctx context.Context, comment *entity.Comment) error {
	return database.DB.WithContext(ctx).Create(comment).Error
}

func (r *CommentRepositoryImpl) Delete(ctx context.Context, id uint) error {
	return database.DB.WithContext(ctx).Delete(&entity.Comment{}, "id = ?", id).Error
}

func (r *CommentRepositoryImpl) Update(ctx context.Context, comment *entity.Comment) error {
	return database.DB.WithContext(ctx).Save(comment).Error
}

func (r *CommentRepositoryImpl) GetByID(ctx context.Context, id string) (*entity.Comment, error) {
	var comment entity.Comment
	result := database.DB.WithContext(ctx).Where("id = ?", id).First(&comment)
	if result.Error != nil {
		return nil, result.Error
	}
	return &comment, nil
}

func (r *CommentRepositoryImpl) GetByPostID(ctx context.Context, postID uint) ([]*entity.Comment, error) {
	var comments []*entity.Comment
	result := database.DB.WithContext(ctx).
		Where("post_id = ? AND parent_id IS NULL", postID).
		Order("created_at DESC").
		Find(&comments)
	if result.Error != nil {
		return nil, result.Error
	}
	return comments, nil
}

func (r *CommentRepositoryImpl) GetRepliesByParentID(ctx context.Context, parentID string) ([]*entity.Comment, error) {
	var comments []*entity.Comment
	result := database.DB.WithContext(ctx).
		Where("parent_id = ?", parentID).
		Order("created_at ASC").
		Find(&comments)
	if result.Error != nil {
		return nil, result.Error
	}
	return comments, nil
}

func NewCommentRepository() CommentRepository {
	return &CommentRepositoryImpl{}
}

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

func (r CommentRepositoryImpl) Create(ctx context.Context, comment *entity.Comment) error {
	return database.DB.WithContext(ctx).Create(comment).Error
}

func (r CommentRepositoryImpl) Delete(ctx context.Context, id string) error {
	return database.DB.WithContext(ctx).Delete(&entity.Comment{}, "id = ?", id).Error
}

func (r CommentRepositoryImpl) Update(ctx context.Context, comment *entity.Comment) error {
	return database.DB.WithContext(ctx).Save(comment).Error
}

func NewCommentRepository() CommentRepository {
	return CommentRepositoryImpl{}
}

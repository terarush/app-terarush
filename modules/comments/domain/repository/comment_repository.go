package repository

import (
	"context"
	"go-modular/modules/comments/domain/entity"
)

type CommentRepository interface {
	Create(ctx context.Context, comment *entity.Comment) error
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, comment *entity.Comment) error
	GetByID(ctx context.Context, id string) (*entity.Comment, error)
	GetByPostID(ctx context.Context, postID uint) ([]*entity.Comment, error)
	GetRepliesByParentID(ctx context.Context, parentID string) ([]*entity.Comment, error)
}

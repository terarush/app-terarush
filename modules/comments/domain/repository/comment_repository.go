package repository

import (
	"context"
	"go-modular/modules/comments/domain/entity"
)

type CommentRepository interface {
	Create(ctx context.Context, comment *entity.Comment) error
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, comment *entity.Comment) error
}

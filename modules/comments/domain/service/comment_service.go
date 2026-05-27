package service

import (
	"context"
	"go-modular/modules/comments/domain/entity"
	"go-modular/modules/comments/domain/repository"
)

var (
	ERR_COMMENT_NOT_FOUND = "comment not found"
)

type CommentService struct {
	commentRepo repository.CommentRepository
}

func NewCommentService(commentRepo repository.CommentRepository) *CommentService {
	return &CommentService{
		commentRepo: commentRepo,
	}
}

func (s *CommentService) CreateComment(ctx context.Context, comment *entity.Comment, userId uint) error {
	comment.UserID = userId
	return s.commentRepo.Create(ctx, comment)
}

func (s *CommentService) DeleteComment(ctx context.Context, id uint) error {
	return s.commentRepo.Delete(ctx, id)
}

func (s *CommentService) UpdateComment(ctx context.Context, comment *entity.Comment) error {
	return s.commentRepo.Update(ctx, comment)
}

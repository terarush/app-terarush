package service

import (
	"context"
	"fmt"
	"go-modular/modules/comments/domain/entity"
	"go-modular/modules/comments/domain/repository"
)

var (
	ERR_COMMENT_NOT_FOUND = "comment not found"
	ERR_BLOG_NOT_FOUND    = "blog not found"
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

func (s *CommentService) DeleteComment(ctx context.Context, id string) error {
	return s.commentRepo.Delete(ctx, id)
}

func (s *CommentService) UpdateComment(ctx context.Context, comment *entity.Comment) error {
	return s.commentRepo.Update(ctx, comment)
}

func (s *CommentService) GetCommentByID(ctx context.Context, commentID string) (*entity.Comment, error) {
	comment, err := s.commentRepo.GetByID(ctx, commentID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", ERR_COMMENT_NOT_FOUND, err)
	}
	return comment, nil
}

func (s *CommentService) GetCommentsByPostID(ctx context.Context, postID uint) ([]*entity.Comment, error) {
	return s.commentRepo.GetByPostID(ctx, postID)
}

func (s *CommentService) GetRepliesByCommentID(ctx context.Context, parentID string) ([]*entity.Comment, error) {
	return s.commentRepo.GetRepliesByParentID(ctx, parentID)
}

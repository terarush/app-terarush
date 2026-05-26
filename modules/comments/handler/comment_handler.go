package handler

import (
	"go-modular/internal/pkg/bus"
	"go-modular/internal/pkg/logger"
	"go-modular/internal/pkg/middleware"
	"go-modular/internal/pkg/utils"
	"go-modular/modules/comments/domain/entity"
	"go-modular/modules/comments/domain/service"
	"go-modular/modules/comments/dto/request"

	"github.com/labstack/echo"
)

type CommentHandler struct {
	commentService service.CommentService
	log            *logger.Logger
	event          *bus.EventBus
	r              utils.Response
}

func NewCommentHandler(log *logger.Logger, event *bus.EventBus, commentService service.CommentService) *CommentHandler {
	return &CommentHandler{
		commentService: commentService,
		log:            log,
		event:          event,
		r:              utils.Response{},
	}
}

func (h *CommentHandler) CreateComment(c echo.Context) error {
	ctx := c.Request().Context()
	userId, err := utils.GetUserIDFromContext(ctx)
	if err != nil {
		h.log.Error("Error getting user ID from context", err)
		h.r.UnauthorizedResponse(c, "Unauthorized")
		return nil
	}

	req := new(request.CreateCommentRequest)

	if err := c.Bind(&req); err != nil {
		h.log.Error("Error binding request", err)
		h.r.BadRequestResponse(c, "Invalid request payload")
		return nil
	}

	comment := entity.NewComment(
		req.Content,
		userId,
		req.PostID,
	)

	if err := h.commentService.CreateComment(ctx, comment, userId); err != nil {
		h.log.Error("Error creating comment", err)
		h.r.InternalServerErrorResponse(c, "Failed to create comment")
		return nil
	}

	h.r.CreatedResponse(c, comment, "Comment created successfully")
	return nil
}

func (h *CommentHandler) DeleteComment(c echo.Context) error {
	ctx := c.Request().Context()
	id := c.Param("id")

	if err := h.commentService.DeleteComment(ctx, id); err != nil {
		h.log.Error("Error deleting comment", err)
		h.r.InternalServerErrorResponse(c, "Failed to delete comment")
		return nil
	}

	h.r.SuccessResponse(c, nil, "Comment deleted successfully")
	return nil
}

func (h *CommentHandler) UpdateComment(c echo.Context) error {
	ctx := c.Request().Context()
	id := c.Param("id")

	req := new(request.UpdateCommentRequest)

	if err := c.Bind(&req); err != nil {
		h.log.Error("Error binding request", err)
		h.r.BadRequestResponse(c, "Invalid request payload")
		return nil
	}

	comment := &entity.Comment{
		ID:      id,
		Content: req.Content,
		PostID:  req.PostID,
	}

	if err := h.commentService.UpdateComment(ctx, comment); err != nil {
		h.log.Error("Error updating comment", err)
		h.r.InternalServerErrorResponse(c, "Failed to update comment")
		return nil
	}

	h.r.SuccessResponse(c, comment, "Comment updated successfully")
	return nil
}

func (h *CommentHandler) RegisterRoutes(e *echo.Echo, basePath string) {
	comment := e.Group(basePath+"/comments", middleware.Auth)
	comment.POST("", h.CreateComment)
	comment.PUT("/:id", h.UpdateComment)
	comment.DELETE("/:id", h.DeleteComment)
}

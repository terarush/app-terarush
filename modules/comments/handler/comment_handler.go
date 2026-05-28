// modules/comments/handler/comment_handler.go

package handler

import (
	"go-modular/internal/pkg/bus"
	"go-modular/internal/pkg/logger"
	"go-modular/internal/pkg/middleware"
	"go-modular/internal/pkg/utils"
	"go-modular/modules/comments/domain/entity"
	"go-modular/modules/comments/domain/service"
	"go-modular/modules/comments/dto/request"
	"strconv"

	"github.com/labstack/echo/v4"
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

// CreateComment creates a new comment
// @Summary Create a new comment
// @Description Add a comment to a blog post
// @Tags Comments
// @Accept json
// @Produce json
// @Param request body request.CreateCommentRequest true "Comment creation request"
// @Success 201 {object} entity.Comment "Comment created successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/comments [post]
// @Security Bearer
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

// DeleteComment deletes a comment
// @Summary Delete a comment
// @Description Remove a comment from a blog post
// @Tags Comments
// @Produce json
// @Param id path int true "Comment ID"
// @Success 200 {object} map[string]interface{} "Comment deleted successfully"
// @Failure 400 {object} map[string]interface{} "Invalid comment ID"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/comments/{id} [delete]
// @Security Bearer
func (h *CommentHandler) DeleteComment(c echo.Context) error {
	ctx := c.Request().Context()
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		h.log.Error("Error parsing comment id", err)
		h.r.BadRequestResponse(c, "Invalid comment id")
		return nil
	}

	if err := h.commentService.DeleteComment(ctx, uint(id)); err != nil {
		h.log.Error("Error deleting comment", err)
		h.r.InternalServerErrorResponse(c, "Failed to delete comment")
		return nil
	}

	h.r.SuccessResponse(c, nil, "Comment deleted successfully")
	return nil
}

// UpdateComment updates a comment
// @Summary Update a comment
// @Description Update an existing comment on a blog post
// @Tags Comments
// @Accept json
// @Produce json
// @Param id path string true "Comment ID"
// @Param request body request.UpdateCommentRequest true "Comment update request"
// @Success 200 {object} entity.Comment "Comment updated successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/comments/{id} [put]
// @Security Bearer
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

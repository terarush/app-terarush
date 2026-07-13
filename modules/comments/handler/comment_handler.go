// modules/comments/handler/comment_handler.go

package handler

import (
	"fmt"
	"go-modular/internal/pkg/bus"
	"go-modular/internal/pkg/logger"
	"go-modular/internal/pkg/middleware"
	"go-modular/internal/pkg/utils"
	"go-modular/modules/comments/domain/entity"
	"go-modular/modules/comments/domain/service"
	"go-modular/modules/comments/dto/request"
	"go-modular/modules/comments/dto/response"
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

// getUserAvatarFromClaims extracts avatar URL from JWT claims
func getUserAvatarFromClaims(c echo.Context) string {
	claims, ok := c.Get("user").(map[string]interface{})
	if !ok {
		return ""
	}
	avatar, _ := claims["avatar"].(string)
	return avatar
}

// getUserIDAndNameFromContext extracts user ID and name from Echo context
func getUserIDAndNameFromContext(c echo.Context) (uint, string, error) {
	claims, ok := c.Get("user").(map[string]interface{})
	if !ok {
		return 0, "", fmt.Errorf("user claims not found in context")
	}

	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return 0, "", fmt.Errorf("user ID not found in claims")
	}
	userID := uint(userIDFloat)

	userName, ok := claims["name"].(string)
	if !ok {
		return 0, "", fmt.Errorf("user name not found in claims")
	}

	return userID, userName, nil
}

// GetCommentsByBlog gets all comments for a blog post
// @Summary Get comments for a blog post
// @Description Retrieve all comments (with nested replies) for a specific blog post
// @Tags Comments
// @Produce json
// @Param post_id query int true "Blog post ID"
// @Success 200 {object} map[string]interface{} "List of comments with replies"
// @Failure 400 {object} map[string]interface{} "Invalid post ID"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/blogs/{post_id}/comments [get]
func (h *CommentHandler) GetCommentsByBlog(c echo.Context) error {
	ctx := c.Request().Context()

	postID, err := strconv.ParseUint(c.Param("post_id"), 10, 32)
	if err != nil {
		h.log.Error("Error parsing post id", err)
		return h.r.BadRequestResponse(c, "Invalid post ID")
	}

	comments, err := h.commentService.GetCommentsByPostID(ctx, uint(postID))
	if err != nil {
		h.log.Error("Error fetching comments", err)
		return h.r.InternalServerErrorResponse(c, "Failed to fetch comments")
	}

	// Fetch replies for each comment
	commentResponses := make([]*response.CommentResponse, 0)
	for _, comment := range comments {
		replies, err := h.commentService.GetRepliesByCommentID(ctx, comment.ID)
		if err != nil {
			h.log.Error("Error fetching replies", err)
			continue
		}

		commentResp := response.NewCommentResponse(
			comment.ID,
			comment.Content,
			comment.UserID,
			comment.PostID,
			comment.UserName,
			comment.CreatedAt,
			comment.UpdatedAt,
			comment.UserAvatar,
		)

		// Convert replies
		commentResp.Replies = make([]*response.CommentResponse, len(replies))
		for i, reply := range replies {
			commentResp.Replies[i] = response.NewReplyResponse(
				reply.ID,
				reply.Content,
				reply.UserID,
				reply.PostID,
				*reply.ParentID,
				reply.UserName,
				reply.CreatedAt,
				reply.UpdatedAt,
				reply.UserAvatar,
			)
		}

		commentResponses = append(commentResponses, commentResp)
	}

	return h.r.SuccessResponse(c, map[string]interface{}{
		"comments": commentResponses,
		"total":    len(commentResponses),
	}, "")
}

// CreateComment creates a new comment
// @Summary Create a new comment
// @Description Add a comment to a blog post
// @Tags Comments
// @Accept json
// @Produce json
// @Param request body request.CreateCommentRequest true "Comment creation request"
// @Success 201 {object} response.CommentResponse "Comment created successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/comments [post]
// @Security Bearer
func (h *CommentHandler) CreateComment(c echo.Context) error {
	ctx := c.Request().Context()

	userID, userName, err := getUserIDAndNameFromContext(c)
	if err != nil {
		h.log.Error("Error getting user info from context", err)
		return h.r.UnauthorizedResponse(c, "Unauthorized")
	}

	req := new(request.CreateCommentRequest)

	if err := c.Bind(&req); err != nil {
		h.log.Error("Error binding request", err)
		return h.r.BadRequestResponse(c, "Invalid request payload")
	}

	comment := entity.NewComment(
		req.Content,
		userID,
		req.PostID,
	)
	comment.UserName = userName
	comment.UserAvatar = getUserAvatarFromClaims(c)

	if err := h.commentService.CreateComment(ctx, comment, userID); err != nil {
		h.log.Error("Error creating comment", err)
		return h.r.InternalServerErrorResponse(c, "Failed to create comment")
	}

	h.event.Publish(bus.Event{Type: "comment.created", Payload: comment})

	commentResp := response.NewCommentResponse(
		comment.ID,
		comment.Content,
		comment.UserID,
		comment.PostID,
		comment.UserName,
		comment.CreatedAt,
		comment.UpdatedAt,
		comment.UserAvatar,
	)

	return h.r.CreatedResponse(c, commentResp, "Comment created successfully")
}

// CreateReply creates a reply to a comment
// @Summary Create a reply to a comment
// @Description Add a reply to an existing comment on a blog post
// @Tags Comments
// @Accept json
// @Produce json
// @Param request body request.CreateReplyRequest true "Reply creation request"
// @Success 201 {object} response.CommentResponse "Reply created successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/comments/reply [post]
// @Security Bearer
func (h *CommentHandler) CreateReply(c echo.Context) error {
	ctx := c.Request().Context()

	userID, userName, err := getUserIDAndNameFromContext(c)
	if err != nil {
		h.log.Error("Error getting user info from context", err)
		return h.r.UnauthorizedResponse(c, "Unauthorized")
	}

	req := new(request.CreateReplyRequest)

	if err := c.Bind(&req); err != nil {
		h.log.Error("Error binding request", err)
		return h.r.BadRequestResponse(c, "Invalid request payload")
	}

	// Validate parent comment exists
	_, err = h.commentService.GetCommentByID(ctx, req.ParentID)
	if err != nil {
		h.log.Error("Parent comment not found", err)
		return h.r.BadRequestResponse(c, "Parent comment not found")
	}

	reply := entity.NewReply(
		req.Content,
		userID,
		req.PostID,
		req.ParentID,
		userName,
	)
	reply.UserAvatar = getUserAvatarFromClaims(c)

	if err := h.commentService.CreateComment(ctx, reply, userID); err != nil {
		h.log.Error("Error creating reply", err)
		return h.r.InternalServerErrorResponse(c, "Failed to create reply")
	}

	h.event.Publish(bus.Event{Type: "comment.reply_created", Payload: reply})

	replyResp := response.NewReplyResponse(
		reply.ID,
		reply.Content,
		reply.UserID,
		reply.PostID,
		*reply.ParentID,
		reply.UserName,
		reply.CreatedAt,
		reply.UpdatedAt,
		reply.UserAvatar,
	)

	return h.r.CreatedResponse(c, replyResp, "Reply created successfully")
}

// DeleteComment deletes a comment
// @Summary Delete a comment
// @Description Remove a comment from a blog post
// @Tags Comments
// @Produce json
// @Param id path string true "Comment ID"
// @Success 200 {object} map[string]interface{} "Comment deleted successfully"
// @Failure 400 {object} map[string]interface{} "Invalid comment ID"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/comments/{id} [delete]
// @Security Bearer
func (h *CommentHandler) DeleteComment(c echo.Context) error {
	ctx := c.Request().Context()
	id := c.Param("id")

	if id == "" {
		h.log.Error("Comment ID is empty")
		return h.r.BadRequestResponse(c, "Invalid comment id")
	}

	// First get the comment to check ownership or admin status
	comment, err := h.commentService.GetCommentByID(ctx, id)
	if err != nil {
		h.log.Error("Error getting comment", err)
		return h.r.BadRequestResponse(c, "Comment not found")
	}

	// Extract user ID for authorization check
	userID, _, err := getUserIDAndNameFromContext(c)
	if err != nil {
		return h.r.UnauthorizedResponse(c, "Unauthorized")
	}

	// Only comment owner or admin can delete
	if comment.UserID != userID {
		return h.r.UnauthorizedResponse(c, "You can only delete your own comments")
	}

	if err := h.commentService.DeleteComment(ctx, id); err != nil {
		h.log.Error("Error deleting comment", err)
		return h.r.InternalServerErrorResponse(c, "Failed to delete comment")
	}

	h.event.Publish(bus.Event{Type: "comment.deleted", Payload: map[string]string{"id": id}})

	return h.r.SuccessResponse(c, nil, "Comment deleted successfully")
}

// UpdateComment updates a comment
// @Summary Update a comment
// @Description Update an existing comment on a blog post
// @Tags Comments
// @Accept json
// @Produce json
// @Param id path string true "Comment ID"
// @Param request body request.UpdateCommentRequest true "Comment update request"
// @Success 200 {object} response.CommentResponse "Comment updated successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/comments/{id} [put]
// @Security Bearer
func (h *CommentHandler) UpdateComment(c echo.Context) error {
	ctx := c.Request().Context()
	id := c.Param("id")

	// Get user info for authorization
	userID, _, err := getUserIDAndNameFromContext(c)
	if err != nil {
		return h.r.UnauthorizedResponse(c, "Unauthorized")
	}

	// Get existing comment
	existingComment, err := h.commentService.GetCommentByID(ctx, id)
	if err != nil {
		h.log.Error("Error getting comment", err)
		return h.r.BadRequestResponse(c, "Comment not found")
	}

	// Check ownership
	if existingComment.UserID != userID {
		return h.r.UnauthorizedResponse(c, "You can only update your own comments")
	}

	req := new(request.UpdateCommentRequest)

	if err := c.Bind(&req); err != nil {
		h.log.Error("Error binding request", err)
		return h.r.BadRequestResponse(c, "Invalid request payload")
	}

	existingComment.Content = req.Content
	existingComment.PostID = req.PostID

	if err := h.commentService.UpdateComment(ctx, existingComment); err != nil {
		h.log.Error("Error updating comment", err)
		return h.r.InternalServerErrorResponse(c, "Failed to update comment")
	}

	h.event.Publish(bus.Event{Type: "comment.updated", Payload: existingComment})

	commentResp := response.NewCommentResponse(
		existingComment.ID,
		existingComment.Content,
		existingComment.UserID,
		existingComment.PostID,
		existingComment.UserName,
		existingComment.CreatedAt,
		existingComment.UpdatedAt,
	)

	return h.r.SuccessResponse(c, commentResp, "Comment updated successfully")
}

func (h *CommentHandler) RegisterRoutes(e *echo.Echo, basePath string) {
	// Public routes - get comments for a blog (no auth required)
	public := e.Group(basePath + "/blogs/:post_id/comments")
	public.GET("", h.GetCommentsByBlog)

	// Authenticated routes
	auth := e.Group(basePath + "/comments", middleware.Auth)
	auth.POST("", h.CreateComment)
	auth.POST("/reply", h.CreateReply)
	auth.PUT("/:id", h.UpdateComment)
	auth.DELETE("/:id", h.DeleteComment)
}

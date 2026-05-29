package request

type CreateCommentRequest struct {
	Content string `json:"content" binding:"required"`
	PostID  uint   `json:"post_id" binding:"required"`
}

type UpdateCommentRequest struct {
	Content string `json:"content" binding:"required"`
	PostID  uint   `json:"post_id" binding:"required"`
}

type CreateReplyRequest struct {
	Content  string `json:"content" binding:"required"`
	PostID   uint   `json:"post_id" binding:"required"`
	ParentID string `json:"parent_id" binding:"required"`
}

package response

type CommentResponse struct {
	ID        string `json:"id"`
	Content   string `json:"content"`
	UserID    uint   `json:"user_id"`
	PostID    uint   `json:"post_id"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

func NewCommentResponse(id, content string, userID, postID uint, createdAt, updatedAt int64) *CommentResponse {
	return &CommentResponse{
		ID:        id,
		Content:   content,
		UserID:    userID,
		PostID:    postID,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}

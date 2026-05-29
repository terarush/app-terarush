package response

type CommentResponse struct {
	ID        string                `json:"id"`
	Content   string                `json:"content"`
	UserID    uint                  `json:"user_id"`
	UserName  string                `json:"user_name"`
	PostID    uint                  `json:"post_id"`
	ParentID  *string               `json:"parent_id"`
	CreatedAt int64                 `json:"created_at"`
	UpdatedAt int64                 `json:"updated_at"`
	Replies   []*CommentResponse    `json:"replies,omitempty"`
}

func NewCommentResponse(id, content string, userID, postID uint, userName string, createdAt, updatedAt int64) *CommentResponse {
	return &CommentResponse{
		ID:        id,
		Content:   content,
		UserID:    userID,
		UserName:  userName,
		PostID:    postID,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}

func NewReplyResponse(id, content string, userID, postID uint, parentID, userName string, createdAt, updatedAt int64) *CommentResponse {
	parentIDPtr := &parentID
	return &CommentResponse{
		ID:        id,
		Content:   content,
		UserID:    userID,
		UserName:  userName,
		PostID:    postID,
		ParentID:  parentIDPtr,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}

package response

type CommentResponse struct {
	ID        string             `json:"id"`
	Content   string             `json:"content"`
	UserID    uint               `json:"user_id"`
	UserName  string             `json:"user_name"`
	UserAvatar string            `json:"user_avatar"`
	PostID    uint               `json:"post_id"`
	ParentID  *string            `json:"parent_id"`
	CreatedAt int64              `json:"created_at"`
	UpdatedAt int64              `json:"updated_at"`
	Replies   []*CommentResponse `json:"replies,omitempty"`
}

func NewCommentResponse(id, content string, userID, postID uint, userName string, createdAt, updatedAt int64, userAvatar ...string) *CommentResponse {
	resp := &CommentResponse{
		ID:        id,
		Content:   content,
		UserID:    userID,
		UserName:  userName,
		PostID:    postID,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
	if len(userAvatar) > 0 {
		resp.UserAvatar = userAvatar[0]
	}
	return resp
}

func NewReplyResponse(id, content string, userID, postID uint, parentID, userName string, createdAt, updatedAt int64, userAvatar ...string) *CommentResponse {
	parentIDPtr := &parentID
	resp := &CommentResponse{
		ID:        id,
		Content:   content,
		UserID:    userID,
		UserName:  userName,
		PostID:    postID,
		ParentID:  parentIDPtr,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
	if len(userAvatar) > 0 {
		resp.UserAvatar = userAvatar[0]
	}
	return resp
}

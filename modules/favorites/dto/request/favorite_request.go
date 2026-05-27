package request

type CreateFavoriteRequest struct {
	UserID uint `json:"user_id"`
	PostID uint `json:"post_id"`
}

package response

import "go-modular/modules/favorites/domain/entity"

type FavoriteResponse struct {
	ID        uint  `json:"id"`
	UserID    uint  `json:"user_id"`
	PostID    uint  `json:"post_id"`
	CreatedAt int64 `json:"created_at"`
}

func NewFavoriteResponse(favorite entity.Favorite) FavoriteResponse {
	return FavoriteResponse{
		ID:        favorite.ID,
		UserID:    favorite.UserID,
		PostID:    favorite.PostID,
		CreatedAt: favorite.CreatedAt,
	}
}

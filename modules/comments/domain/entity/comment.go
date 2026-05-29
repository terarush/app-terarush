package entity

import "time"

type Comment struct {
	ID          string `gorm:"primaryKey" json:"id"`
	Content     string `gorm:"type:text" json:"content"`
	UserID      uint   `gorm:"index" json:"user_id"`
	PostID      uint   `gorm:"index;constraint:OnDelete:CASCADE" json:"post_id"` // Foreign key dengan cascade delete
	ParentID    *string `gorm:"index" json:"parent_id"`        // For reply: references another comment ID
	UserName    string  `gorm:"type:varchar(100)" json:"user_name"` // Store commenter's name
	CreatedAt   int64   `json:"created_at"`
	UpdatedAt   int64   `json:"updated_at"`
	DeletedAt   *int64  `json:"deleted_at"`
}

func (*Comment) TableName() string {
	return "comments"
}

func NewComment(content string, userID, postID uint) *Comment {
	timestamp := time.Now().Unix()
	return &Comment{
		Content:   content,
		UserID:    userID,
		PostID:    postID,
		CreatedAt: timestamp,
		UpdatedAt: timestamp,
	}
}

func NewReply(content string, userID, postID uint, parentID string, userName string) *Comment {
	timestamp := time.Now().Unix()
	return &Comment{
		Content:   content,
		UserID:    userID,
		PostID:    postID,
		ParentID:  &parentID,
		UserName:  userName,
		CreatedAt: timestamp,
		UpdatedAt: timestamp,
	}
}

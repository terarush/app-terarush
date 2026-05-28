package entity

import "time"

type Asset struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	URL       string    `json:"url"`
	Path      string    `json:"path,omitempty"`
	FileName  string    `json:"file_name"`
	FileSize  int64     `json:"file_size"`
	MimeType  string    `json:"mime_type"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (Asset) TableName() string {
	return "assets"
}

func NewAsset(url string) *Asset {
	return &Asset{
		URL:       url,
		CreatedAt: time.Now(),
	}
}

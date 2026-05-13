package models

import "time"

type PostStatus string

const (
	PostStatusDraft     PostStatus = "draft"
	PostStatusPublished PostStatus = "published"
)

// Post — новость / запись блога института. Контент на 3 языках.
type Post struct {
	ID          int64      `gorm:"primaryKey" json:"id"`
	Slug        string     `gorm:"size:255;uniqueIndex;not null" json:"slug"`
	CoverPath   string     `gorm:"type:text" json:"cover_path,omitempty"`
	Author      string     `gorm:"size:255" json:"author,omitempty"`
	PublishedAt *time.Time `gorm:"index" json:"published_at,omitempty"`
	Status      PostStatus `gorm:"type:post_status;not null;default:'draft';index" json:"status"`
	TitleTM     string     `gorm:"size:512;not null" json:"title_tm"`
	TitleRU     string     `gorm:"size:512;not null" json:"title_ru"`
	TitleEN     string     `gorm:"size:512;not null" json:"title_en"`
	ContentTM   string     `gorm:"type:text;not null;default:''" json:"content_tm"`
	ContentRU   string     `gorm:"type:text;not null;default:''" json:"content_ru"`
	ContentEN   string     `gorm:"type:text;not null;default:''" json:"content_en"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func (Post) TableName() string { return "posts" }

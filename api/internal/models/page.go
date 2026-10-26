package models

import "time"

// Page — статичная страница CMS ("О библиотеке", "Контакты", "Правила").
// Контент через TipTap WYSIWYG, на 3 языках.
// Опционально добавляется в меню сайта (флаг InMenu) как самостоятельный пункт.
type Page struct {
	ID        int64     `gorm:"primaryKey" json:"id"`
	Slug      string    `gorm:"size:255;uniqueIndex;not null" json:"slug"`
	InMenu    bool      `gorm:"not null;default:false;index" json:"in_menu"`
	MenuOrder int       `gorm:"not null;default:0" json:"menu_order"`
	TitleTM   string    `gorm:"size:255;not null" json:"title_tm"`
	TitleRU   string    `gorm:"size:255;not null" json:"title_ru"`
	TitleEN   string    `gorm:"size:255;not null" json:"title_en"`
	ContentTM string    `gorm:"type:text;not null;default:''" json:"content_tm"`
	ContentRU string    `gorm:"type:text;not null;default:''" json:"content_ru"`
	ContentEN string    `gorm:"type:text;not null;default:''" json:"content_en"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Page) TableName() string { return "pages" }

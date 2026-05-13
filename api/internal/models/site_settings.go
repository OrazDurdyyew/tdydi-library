package models

import (
	"time"

	"gorm.io/datatypes"
)

// SiteSettings — глобальные настройки сайта (singleton, всегда одна строка с id=1).
// Хранит брендинг (логотип, фавикон, цветовые токены), контакты, название института.
// Цвета хранятся в формате "R G B" (например, "16 117 235") для совместимости
// с rgb(var(--primary) / <alpha-value>) в Tailwind / shadcn-vue.
type SiteSettings struct {
	ID              int16          `gorm:"primaryKey;default:1" json:"id"`
	NameTM          string         `gorm:"size:255;not null;default:''" json:"name_tm"`
	NameRU          string         `gorm:"size:255;not null;default:''" json:"name_ru"`
	NameEN          string         `gorm:"size:255;not null;default:''" json:"name_en"`
	LogoLightPath   string         `gorm:"type:text" json:"logo_light_path,omitempty"`
	LogoDarkPath    string         `gorm:"type:text" json:"logo_dark_path,omitempty"`
	FaviconPath     string         `gorm:"type:text" json:"favicon_path,omitempty"`
	PrimaryColor    string         `gorm:"size:32;not null" json:"primary_color"`
	SecondaryColor  string         `gorm:"size:32;not null" json:"secondary_color"`
	AccentColor     string         `gorm:"size:32;not null" json:"accent_color"`
	BackgroundColor string         `gorm:"size:32;not null" json:"background_color"`
	ForegroundColor string         `gorm:"size:32;not null" json:"foreground_color"`
	MutedColor      string         `gorm:"size:32;not null" json:"muted_color"`
	BorderColor     string         `gorm:"size:32;not null" json:"border_color"`
	RingColor       string         `gorm:"size:32;not null" json:"ring_color"`
	FontFamily      string         `gorm:"size:255;not null;default:'Inter'" json:"font_family"`
	ContactAddress  string         `gorm:"type:text" json:"contact_address,omitempty"`
	ContactPhone    string         `gorm:"size:64" json:"contact_phone,omitempty"`
	ContactEmail    string         `gorm:"size:255" json:"contact_email,omitempty"`
	SocialLinks     datatypes.JSON `gorm:"type:jsonb;not null;default:'{}'" json:"social_links"`
	FooterTM        string         `gorm:"type:text;not null;default:''" json:"footer_tm"`
	FooterRU        string         `gorm:"type:text;not null;default:''" json:"footer_ru"`
	FooterEN        string         `gorm:"type:text;not null;default:''" json:"footer_en"`
	UpdatedAt       time.Time      `json:"updated_at"`
}

func (SiteSettings) TableName() string { return "site_settings" }

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
	// Палитра тёмной темы (светлая — поля выше).
	DarkPrimaryColor    string `gorm:"size:32;not null;default:'122 162 201'" json:"dark_primary_color"`
	DarkSecondaryColor  string `gorm:"size:32;not null;default:'96 141 184'" json:"dark_secondary_color"`
	DarkAccentColor     string `gorm:"size:32;not null;default:'200 169 107'" json:"dark_accent_color"`
	DarkBackgroundColor string `gorm:"size:32;not null;default:'15 23 42'" json:"dark_background_color"`
	DarkForegroundColor string `gorm:"size:32;not null;default:'249 250 251'" json:"dark_foreground_color"`
	DarkMutedColor      string `gorm:"size:32;not null;default:'30 41 59'" json:"dark_muted_color"`
	DarkBorderColor     string `gorm:"size:32;not null;default:'31 41 55'" json:"dark_border_color"`
	DarkRingColor       string `gorm:"size:32;not null;default:'122 162 201'" json:"dark_ring_color"`
	FontFamily      string         `gorm:"size:255;not null;default:'Inter'" json:"font_family"`
	FontFilePath    string         `gorm:"type:text" json:"font_file_path,omitempty"`
	FontURL         string         `gorm:"type:text" json:"font_url,omitempty"`
	FontUploadedAt  time.Time      `json:"font_uploaded_at,omitempty"`
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

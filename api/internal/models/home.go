package models

import "time"

// HomeSlide — слайд баннера на главной странице.
type HomeSlide struct {
	ID         int64     `gorm:"primaryKey" json:"id"`
	ImagePath  string    `gorm:"type:text;not null;default:''" json:"image_path"`
	SortOrder  int       `gorm:"not null;default:0" json:"sort_order"`
	IsActive   bool      `gorm:"not null;default:true" json:"is_active"`
	LinkURL    string    `gorm:"type:text" json:"link_url,omitempty"`
	TitleTM    string    `gorm:"size:255;not null;default:''" json:"title_tm"`
	TitleRU    string    `gorm:"size:255;not null;default:''" json:"title_ru"`
	TitleEN    string    `gorm:"size:255;not null;default:''" json:"title_en"`
	SubtitleTM string    `gorm:"type:text;not null;default:''" json:"subtitle_tm"`
	SubtitleRU string    `gorm:"type:text;not null;default:''" json:"subtitle_ru"`
	SubtitleEN string    `gorm:"type:text;not null;default:''" json:"subtitle_en"`
	CreatedAt  time.Time `json:"created_at"`
}

// HomeSection — какой раздел (или подраздел) показывать на главной как блок.
type HomeSection struct {
	ID        int64 `gorm:"primaryKey" json:"id"`
	SectionID int64 `gorm:"not null;uniqueIndex" json:"section_id"`
	SortOrder int   `gorm:"not null;default:0" json:"sort_order"`
	IsActive  bool  `gorm:"not null;default:true" json:"is_active"`

	Section *Section `gorm:"foreignKey:SectionID" json:"section,omitempty"`
}

// HomeCollection — какую подборку показывать на главной.
type HomeCollection struct {
	ID           int64 `gorm:"primaryKey" json:"id"`
	CollectionID int64 `gorm:"not null;uniqueIndex" json:"collection_id"`
	SortOrder    int   `gorm:"not null;default:0" json:"sort_order"`
	IsActive     bool  `gorm:"not null;default:true" json:"is_active"`

	Collection *Collection `gorm:"foreignKey:CollectionID" json:"collection,omitempty"`
}

func (HomeSlide) TableName() string      { return "home_slides" }
func (HomeSection) TableName() string    { return "home_sections" }
func (HomeCollection) TableName() string { return "home_collections" }

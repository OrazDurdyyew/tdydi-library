package models

import "time"

// Section — раздел или подраздел. 2-уровневая иерархия:
// parent_id IS NULL — это корневой раздел (отображается в главном меню).
// parent_id IS NOT NULL — это подраздел. Подразделы не могут иметь своих подразделов
// (ограничение enforced триггером sections_max_depth_trigger в БД).
type Section struct {
	ID        int64     `gorm:"primaryKey" json:"id"`
	Slug      string    `gorm:"size:255;uniqueIndex;not null" json:"slug"`
	ParentID  *int64    `gorm:"index" json:"parent_id"`
	SortOrder int       `gorm:"not null;default:0" json:"sort_order"`
	ImagePath string    `gorm:"type:text" json:"image_path,omitempty"`
	NameTM    string    `gorm:"size:255;not null" json:"name_tm"`
	NameRU    string    `gorm:"size:255;not null" json:"name_ru"`
	NameEN    string    `gorm:"size:255;not null" json:"name_en"`
	DescTM    string    `gorm:"type:text;column:description_tm;not null;default:''" json:"description_tm"`
	DescRU    string    `gorm:"type:text;column:description_ru;not null;default:''" json:"description_ru"`
	DescEN    string    `gorm:"type:text;column:description_en;not null;default:''" json:"description_en"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Parent   *Section  `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Children []Section `gorm:"foreignKey:ParentID" json:"children,omitempty"`

	// BookCount — не хранится в БД, заполняется хендлерами для админки и карточек подразделов.
	BookCount int64 `gorm:"-" json:"book_count"`
}

func (Section) TableName() string { return "sections" }

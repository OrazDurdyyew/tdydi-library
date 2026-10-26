package models

import "time"

// Collection — виртуальная подборка (например, "К экзамену по микроэкономике").
// Создаётся админом, может содержать произвольные книги, опционально отображается на главной.
type Collection struct {
	ID         int64     `gorm:"primaryKey" json:"id"`
	Slug       string    `gorm:"size:255;uniqueIndex;not null" json:"slug"`
	CoverPath  string    `gorm:"type:text" json:"cover_path,omitempty"`
	SortOrder  int       `gorm:"not null;default:0" json:"sort_order"`
	ShowOnHome bool      `gorm:"not null;default:false;index" json:"show_on_home"`
	NameTM     string    `gorm:"size:255;not null" json:"name_tm"`
	NameRU     string    `gorm:"size:255;not null" json:"name_ru"`
	NameEN     string    `gorm:"size:255;not null" json:"name_en"`
	DescTM     string    `gorm:"type:text;column:description_tm;not null;default:''" json:"description_tm"`
	DescRU     string    `gorm:"type:text;column:description_ru;not null;default:''" json:"description_ru"`
	DescEN     string    `gorm:"type:text;column:description_en;not null;default:''" json:"description_en"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`

	Books []Book `gorm:"many2many:collection_books;joinForeignKey:collection_id;joinReferences:book_id" json:"books,omitempty"`
}

type CollectionBook struct {
	CollectionID int64 `gorm:"primaryKey" json:"collection_id"`
	BookID       int64 `gorm:"primaryKey" json:"book_id"`
	SortOrder    int   `gorm:"not null;default:0" json:"sort_order"`
}

func (Collection) TableName() string     { return "collections" }
func (CollectionBook) TableName() string { return "collection_books" }

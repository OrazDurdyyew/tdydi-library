package models

import "time"

// Book — основная сущность каталога. Метаданные на одном языке (поле Language).
// Файл хранится в файловой системе (FilePath), в БД только путь.
// Книга связана с разделами через M:N — может находиться в нескольких разделах/подразделах.
type Book struct {
	ID               int64     `gorm:"primaryKey" json:"id"`
	Slug             string    `gorm:"size:255;uniqueIndex;not null" json:"slug"`
	ISBN             string    `gorm:"size:32" json:"isbn,omitempty"`
	Year             int16     `gorm:"index" json:"year,omitempty"`
	Pages            int       `json:"pages,omitempty"`
	Language         string    `gorm:"size:8;not null;default:'tm';index" json:"language"`
	Title            string    `gorm:"size:512;not null" json:"title"`
	Annotation       string    `gorm:"type:text;not null;default:''" json:"annotation"`
	Authors          string    `gorm:"type:text;not null;default:''" json:"authors"`
	Publisher        string    `gorm:"size:255" json:"publisher,omitempty"`
	FilePath         string    `gorm:"type:text" json:"-"`
	FileFormat       string    `gorm:"size:16" json:"file_format,omitempty"`
	CoverPath        string    `gorm:"type:text" json:"cover_path,omitempty"`
	AllowDownload    bool      `gorm:"not null;default:true" json:"allow_download"`
	IsPublished      bool      `gorm:"not null;default:false;index" json:"is_published"`
	IsInstituteBook  bool      `gorm:"not null;default:false;index" json:"is_institute_book"`
	ViewsCount       int64     `gorm:"not null;default:0" json:"views_count"`
	DownloadsCount   int64     `gorm:"not null;default:0" json:"downloads_count"`
	CreatedByAdminID *int64    `json:"created_by_admin_id,omitempty"`
	CreatedAt        time.Time `gorm:"index" json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`

	Sections []Section `gorm:"many2many:book_sections;joinForeignKey:book_id;joinReferences:section_id" json:"sections,omitempty"`
	Tags     []Tag     `gorm:"many2many:book_tags;joinForeignKey:book_id;joinReferences:tag_id" json:"tags,omitempty"`
}

type Tag struct {
	ID        int64     `gorm:"primaryKey" json:"id"`
	NameTM    string    `gorm:"size:128;not null" json:"name_tm"`
	NameRU    string    `gorm:"size:128;not null" json:"name_ru"`
	NameEN    string    `gorm:"size:128;not null" json:"name_en"`
	CreatedAt time.Time `json:"created_at"`
}

type BookSection struct {
	BookID    int64 `gorm:"primaryKey" json:"book_id"`
	SectionID int64 `gorm:"primaryKey" json:"section_id"`
}

type BookTag struct {
	BookID int64 `gorm:"primaryKey" json:"book_id"`
	TagID  int64 `gorm:"primaryKey" json:"tag_id"`
}

func (Book) TableName() string        { return "books" }
func (Tag) TableName() string         { return "tags" }
func (BookSection) TableName() string { return "book_sections" }
func (BookTag) TableName() string     { return "book_tags" }

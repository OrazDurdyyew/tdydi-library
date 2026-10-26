package models

import (
	"time"

	"gorm.io/gorm"
)

// Department — кафедра. Группирует преподавателей.
type Department struct {
	ID        int64     `gorm:"primaryKey" json:"id"`
	Slug      string    `gorm:"size:255;uniqueIndex;not null" json:"slug"`
	NameTM    string    `gorm:"size:255;not null" json:"name_tm"`
	NameRU    string    `gorm:"size:255;not null" json:"name_ru"`
	NameEN    string    `gorm:"size:255;not null" json:"name_en"`
	DescTM    string    `gorm:"type:text;column:description_tm;not null;default:''" json:"description_tm"`
	DescRU    string    `gorm:"type:text;column:description_ru;not null;default:''" json:"description_ru"`
	DescEN    string    `gorm:"type:text;column:description_en;not null;default:''" json:"description_en"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Teacher — преподаватель. Имя на одном языке, три отдельных поля.
type Teacher struct {
	ID           int64     `gorm:"primaryKey" json:"id"`
	DepartmentID int64     `gorm:"not null;index" json:"department_id"`
	Slug         string    `gorm:"size:255;uniqueIndex;not null" json:"slug"`
	FirstName    string    `gorm:"size:255;not null" json:"first_name"`
	LastName     string    `gorm:"size:255;not null" json:"last_name"`
	MiddleName   string    `gorm:"size:255;not null;default:''" json:"middle_name"`
	PhotoPath    string    `gorm:"type:text" json:"photo_path,omitempty"`
	Bio          string    `gorm:"type:text;not null;default:''" json:"bio"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	Department *Department `gorm:"foreignKey:DepartmentID" json:"department,omitempty"`
}

// PublicationCategory — тип публикации (газета, журнал, конференция). Динамическая таблица.
type PublicationCategory struct {
	ID        int64     `gorm:"primaryKey" json:"id"`
	Slug      string    `gorm:"size:255;uniqueIndex;not null" json:"slug"`
	NameTM    string    `gorm:"size:255;not null" json:"name_tm"`
	NameRU    string    `gorm:"size:255;not null" json:"name_ru"`
	NameEN    string    `gorm:"size:255;not null" json:"name_en"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Publication — статья преподавателя. Заголовок на 3 языках (обязателен только name_tm),
// файл хранится в файловой системе (FilePath), в БД только путь.
type Publication struct {
	ID                int64      `gorm:"primaryKey" json:"id"`
	Slug              string     `gorm:"size:255;uniqueIndex;not null" json:"slug"`
	TeacherID         int64      `gorm:"not null;index" json:"teacher_id"`
	CategoryID        int64      `gorm:"not null;index" json:"category_id"`
	PeriodicalIssueID *int64     `gorm:"index" json:"periodical_issue_id,omitempty"`
	NameTM            string     `gorm:"size:512;not null" json:"name_tm"`
	NameRU            string     `gorm:"size:512;not null;default:''" json:"name_ru"`
	NameEN            string     `gorm:"size:512;not null;default:''" json:"name_en"`
	ArticleDate       *time.Time `gorm:"type:date" json:"article_date,omitempty"`
	PageCount         int        `json:"page_count,omitempty"`
	FilePath          string     `gorm:"type:text" json:"-"`
	CoverPath         string     `gorm:"type:text" json:"cover_path,omitempty"`
	HasFile           bool       `gorm:"-" json:"has_file"`
	IsPublished       bool       `gorm:"not null;default:true;index" json:"is_published"`
	ViewsCount        int64      `gorm:"not null;default:0" json:"views_count"`
	DownloadsCount    int64      `gorm:"not null;default:0" json:"downloads_count"`
	CreatedAt         time.Time  `gorm:"index" json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`

	Teacher         *Teacher             `gorm:"foreignKey:TeacherID" json:"teacher,omitempty"`
	Category        *PublicationCategory `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	PeriodicalIssue *PeriodicalIssue     `gorm:"foreignKey:PeriodicalIssueID" json:"periodical_issue,omitempty"`
}

func (Department) TableName() string          { return "departments" }
func (Teacher) TableName() string             { return "teachers" }
func (PublicationCategory) TableName() string { return "publication_categories" }
func (Publication) TableName() string         { return "publications" }

// AfterFind заполняет вычисляемое HasFile: FilePath скрыт из JSON (json:"-"),
// поэтому фронту нужен булев индикатор наличия загруженного файла.
func (p *Publication) AfterFind(*gorm.DB) error {
	p.HasFile = p.FilePath != ""
	return nil
}

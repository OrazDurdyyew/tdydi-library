package models

import "time"

type PeriodicalType string

const (
	PeriodicalNewspaper PeriodicalType = "newspaper"
	PeriodicalMagazine  PeriodicalType = "magazine"
)

// Periodical — периодическое издание (газета или журнал).
// Объединяет оба типа: различается полем Type.
type Periodical struct {
	ID          int64          `gorm:"primaryKey" json:"id"`
	Slug        string         `gorm:"size:255;uniqueIndex;not null" json:"slug"`
	Type        PeriodicalType `gorm:"type:periodical_type;not null;index" json:"type"`
	CoverPath   string         `gorm:"type:text" json:"cover_path,omitempty"`
	Periodicity string         `gorm:"size:64" json:"periodicity,omitempty"`
	NameTM      string         `gorm:"size:255;not null" json:"name_tm"`
	NameRU      string         `gorm:"size:255;not null" json:"name_ru"`
	NameEN      string         `gorm:"size:255;not null" json:"name_en"`
	DescTM      string         `gorm:"type:text;column:description_tm;not null;default:''" json:"description_tm"`
	DescRU      string         `gorm:"type:text;column:description_ru;not null;default:''" json:"description_ru"`
	DescEN      string         `gorm:"type:text;column:description_en;not null;default:''" json:"description_en"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`

	Issues []PeriodicalIssue `gorm:"foreignKey:PeriodicalID" json:"issues,omitempty"`
}

// PeriodicalIssue — конкретный выпуск (номер) издания.
type PeriodicalIssue struct {
	ID             int64     `gorm:"primaryKey" json:"id"`
	PeriodicalID   int64     `gorm:"not null;index" json:"periodical_id"`
	IssueNumber    string    `gorm:"size:64" json:"issue_number,omitempty"`
	IssueDate      time.Time `gorm:"type:date;not null;index" json:"issue_date"`
	FilePath       string    `gorm:"type:text" json:"-"`
	CoverPath      string    `gorm:"type:text" json:"cover_path,omitempty"`
	AnnotationTM   string    `gorm:"type:text;column:annotation_tm;not null;default:''" json:"annotation_tm"`
	AnnotationRU   string    `gorm:"type:text;column:annotation_ru;not null;default:''" json:"annotation_ru"`
	AnnotationEN   string    `gorm:"type:text;column:annotation_en;not null;default:''" json:"annotation_en"`
	ViewsCount     int64     `gorm:"not null;default:0" json:"views_count"`
	DownloadsCount int64     `gorm:"not null;default:0" json:"downloads_count"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func (Periodical) TableName() string      { return "periodicals" }
func (PeriodicalIssue) TableName() string { return "periodical_issues" }

package models

import "time"

type FAQCategory struct {
	ID        int64     `gorm:"primaryKey" json:"id"`
	SortOrder int       `gorm:"not null;default:0" json:"sort_order"`
	NameTM    string    `gorm:"size:255;not null" json:"name_tm"`
	NameRU    string    `gorm:"size:255;not null" json:"name_ru"`
	NameEN    string    `gorm:"size:255;not null" json:"name_en"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Items []FAQItem `gorm:"foreignKey:CategoryID" json:"items,omitempty"`
}

type FAQItem struct {
	ID         int64     `gorm:"primaryKey" json:"id"`
	CategoryID int64     `gorm:"not null;index" json:"category_id"`
	SortOrder  int       `gorm:"not null;default:0" json:"sort_order"`
	QuestionTM string    `gorm:"type:text;not null" json:"question_tm"`
	QuestionRU string    `gorm:"type:text;not null" json:"question_ru"`
	QuestionEN string    `gorm:"type:text;not null" json:"question_en"`
	AnswerTM   string    `gorm:"type:text;not null" json:"answer_tm"`
	AnswerRU   string    `gorm:"type:text;not null" json:"answer_ru"`
	AnswerEN   string    `gorm:"type:text;not null" json:"answer_en"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (FAQCategory) TableName() string { return "faq_categories" }
func (FAQItem) TableName() string     { return "faq_items" }

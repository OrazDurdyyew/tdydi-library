package models

import "time"

type ViewEntityType string

const (
	ViewEntityBook            ViewEntityType = "book"
	ViewEntityPeriodicalIssue ViewEntityType = "periodical_issue"
	ViewEntitySection         ViewEntityType = "section"
	ViewEntityPost            ViewEntityType = "post"
)

// ViewLog — лог просмотров для агрегированной статистики.
// IPHash — sha256(ip + daily_salt), позволяет дедуплицировать просмотры в пределах дня
// без хранения чистых IP (для приватности).
type ViewLog struct {
	ID         int64          `gorm:"primaryKey" json:"id"`
	EntityType ViewEntityType `gorm:"type:view_entity_type;not null" json:"entity_type"`
	EntityID   int64          `gorm:"not null" json:"entity_id"`
	IPHash     string         `gorm:"size:64" json:"-"`
	CreatedAt  time.Time      `json:"created_at"`
}

func (ViewLog) TableName() string { return "view_logs" }

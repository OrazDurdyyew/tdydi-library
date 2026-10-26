package models

import "time"

type AdminRole string

const (
	RoleSuperadmin         AdminRole = "superadmin"
	RoleContentManager     AdminRole = "content_manager"
	RolePublicationManager AdminRole = "publication_manager"
)

type Admin struct {
	ID           int64     `gorm:"primaryKey" json:"id"`
	Username     string    `gorm:"size:64;uniqueIndex;not null" json:"username"`
	PasswordHash string    `gorm:"not null" json:"-"`
	Email        string    `gorm:"size:255" json:"email,omitempty"`
	Role         AdminRole `gorm:"type:admin_role;not null;default:'content_manager'" json:"role"`
	AllSections  bool      `gorm:"not null;default:false" json:"all_sections"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	SectionPermissions []AdminSectionPermission `gorm:"foreignKey:AdminID" json:"section_permissions,omitempty"`
}

type AdminSectionPermission struct {
	AdminID   int64 `gorm:"primaryKey" json:"admin_id"`
	SectionID int64 `gorm:"primaryKey" json:"section_id"`
}

func (Admin) TableName() string                  { return "admins" }
func (AdminSectionPermission) TableName() string { return "admin_section_permissions" }

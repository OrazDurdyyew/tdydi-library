package services

import (
	"gorm.io/gorm"

	"library/internal/models"
)

type PermissionsService struct {
	db *gorm.DB
}

func NewPermissionsService(db *gorm.DB) *PermissionsService {
	return &PermissionsService{db: db}
}

// CanManageSections возвращает true, если админ имеет право управлять книгами
// во всех указанных разделах. Суперадмин и admin.AllSections=true — всегда true.
// Если sectionIDs пустой — false (книга обязана быть хотя бы в одном разделе).
func (s *PermissionsService) CanManageSections(adminID int64, sectionIDs []int64) (bool, error) {
	var admin models.Admin
	if err := s.db.First(&admin, adminID).Error; err != nil {
		return false, err
	}
	if admin.Role == models.RoleSuperadmin || admin.AllSections {
		return true, nil
	}
	if len(sectionIDs) == 0 {
		return false, nil
	}

	var count int64
	if err := s.db.Model(&models.AdminSectionPermission{}).
		Where("admin_id = ? AND section_id IN ?", adminID, sectionIDs).
		Count(&count).Error; err != nil {
		return false, err
	}
	return int(count) == len(sectionIDs), nil
}

// AllowedSectionIDs возвращает список разрешённых разделов админа.
// Для суперадмина/all_sections — возвращает nil (что означает "все разделы").
func (s *PermissionsService) AllowedSectionIDs(adminID int64) ([]int64, error) {
	var admin models.Admin
	if err := s.db.First(&admin, adminID).Error; err != nil {
		return nil, err
	}
	if admin.Role == models.RoleSuperadmin || admin.AllSections {
		return nil, nil
	}
	var ids []int64
	err := s.db.Model(&models.AdminSectionPermission{}).
		Where("admin_id = ?", adminID).
		Pluck("section_id", &ids).Error
	return ids, err
}

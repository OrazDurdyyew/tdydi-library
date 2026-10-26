package handlers

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"library/internal/models"
)

type PublicationStatsHandler struct {
	db *gorm.DB
}

func NewPublicationStatsHandler(db *gorm.DB) *PublicationStatsHandler {
	return &PublicationStatsHandler{db: db}
}

type TopTeacherStat struct {
	ID               int64  `json:"id"`
	Name             string `json:"name"`
	PublicationCount int64  `json:"publication_count"`
}

type PerDepartmentStat struct {
	Department string `json:"department"`
	Count      int64  `json:"count"`
}

// Overview — GET /api/admin/stats/publications. Сводка по публикациям для дашборда.
func (h *PublicationStatsHandler) Overview(c *fiber.Ctx) error {
	var totalPublications, totalTeachers, totalDepartments int64
	h.db.Model(&models.Publication{}).Count(&totalPublications)
	h.db.Model(&models.Teacher{}).Count(&totalTeachers)
	h.db.Model(&models.Department{}).Count(&totalDepartments)

	topTeachers := []TopTeacherStat{}
	h.db.Table("teachers").
		Select("teachers.id AS id, TRIM(teachers.last_name || ' ' || teachers.first_name) AS name, COUNT(publications.id) AS publication_count").
		Joins("LEFT JOIN publications ON publications.teacher_id = teachers.id").
		Group("teachers.id").
		Order("publication_count DESC").
		Limit(10).
		Scan(&topTeachers)

	perDepartment := []PerDepartmentStat{}
	h.db.Table("departments").
		Select("departments.name_tm AS department, COUNT(publications.id) AS count").
		Joins("LEFT JOIN teachers ON teachers.department_id = departments.id").
		Joins("LEFT JOIN publications ON publications.teacher_id = teachers.id").
		Group("departments.id").
		Order("count DESC").
		Scan(&perDepartment)

	return c.JSON(fiber.Map{
		"total_publications": totalPublications,
		"total_teachers":     totalTeachers,
		"total_departments":  totalDepartments,
		"top_teachers":       topTeachers,
		"per_department":     perDepartment,
	})
}

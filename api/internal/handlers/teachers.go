package handlers

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"library/internal/models"
	"library/internal/services"
)

type TeachersHandler struct {
	db    *gorm.DB
	files *services.FileStorage
}

func NewTeachersHandler(db *gorm.DB, files *services.FileStorage) *TeachersHandler {
	return &TeachersHandler{db: db, files: files}
}

type TeacherDTO struct {
	DepartmentID int64  `json:"department_id"`
	Slug         string `json:"slug"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	MiddleName   string `json:"middle_name"`
	Bio          string `json:"bio"`
}

func (h *TeachersHandler) List(c *fiber.Ctx) error {
	q := h.db.Model(&models.Teacher{}).Preload("Department").Order("last_name ASC, first_name ASC")
	if dep := c.Query("department_id"); dep != "" {
		if id, err := strconv.ParseInt(dep, 10, 64); err == nil {
			q = q.Where("department_id = ?", id)
		}
	}

	var total int64
	q.Count(&total)

	// Пагинация применяется только при явном ?page — без него отдаём всё (источник для дропдаунов).
	page, limit := 1, int(total)
	if c.Query("page") != "" {
		page, _ = strconv.Atoi(c.Query("page", "1"))
		limit, _ = strconv.Atoi(c.Query("limit", "20"))
		if page < 1 {
			page = 1
		}
		if limit < 1 || limit > 100 {
			limit = 20
		}
		q = q.Offset((page - 1) * limit).Limit(limit)
	}

	var items []models.Teacher
	q.Find(&items)
	return c.JSON(fiber.Map{"data": items, "meta": fiber.Map{"total": total, "page": page, "limit": limit}})
}

func (h *TeachersHandler) Get(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid id"})
	}
	var t models.Teacher
	if err := h.db.Preload("Department").First(&t, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "not found"})
	}
	return c.JSON(t)
}

func (h *TeachersHandler) GetBySlug(c *fiber.Ctx) error {
	var t models.Teacher
	if err := h.db.Preload("Department").Where("slug = ?", c.Params("slug")).First(&t).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "not found"})
	}
	return c.JSON(t)
}

// uniqueTeacherSlug гарантирует уникальность slug, добавляя числовой суффикс при коллизии.
func (h *TeachersHandler) uniqueTeacherSlug(base string, excludeID int64) string {
	slug := base
	for i := 2; ; i++ {
		var count int64
		q := h.db.Model(&models.Teacher{}).Where("slug = ?", slug)
		if excludeID > 0 {
			q = q.Where("id <> ?", excludeID)
		}
		q.Count(&count)
		if count == 0 {
			return slug
		}
		slug = fmt.Sprintf("%s-%d", base, i)
	}
}

func (h *TeachersHandler) Create(c *fiber.Ctx) error {
	var dto TeacherDTO
	if err := c.BodyParser(&dto); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	if dto.FirstName == "" || dto.LastName == "" {
		return c.Status(400).JSON(fiber.Map{"error": "first_name and last_name required"})
	}
	var dep models.Department
	if err := h.db.First(&dep, dto.DepartmentID).Error; err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "department not found"})
	}
	slug := dto.Slug
	if slug == "" {
		slug = dep.Slug + "-" + slugify(dto.FirstName+" "+dto.LastName)
	}
	slug = h.uniqueTeacherSlug(slug, 0)

	t := models.Teacher{
		DepartmentID: dto.DepartmentID,
		Slug:         slug,
		FirstName:    dto.FirstName,
		LastName:     dto.LastName,
		MiddleName:   dto.MiddleName,
		Bio:          dto.Bio,
	}
	if err := h.db.Create(&t).Error; err != nil {
		return c.Status(409).JSON(fiber.Map{"error": err.Error()})
	}
	h.db.Preload("Department").First(&t, t.ID)
	return c.Status(201).JSON(t)
}

func (h *TeachersHandler) Update(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid id"})
	}
	var t models.Teacher
	if err := h.db.First(&t, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "not found"})
	}
	var dto TeacherDTO
	if err := c.BodyParser(&dto); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	if dto.DepartmentID != 0 {
		var dep models.Department
		if err := h.db.First(&dep, dto.DepartmentID).Error; err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "department not found"})
		}
		t.DepartmentID = dto.DepartmentID
	}
	if dto.Slug != "" {
		t.Slug = h.uniqueTeacherSlug(dto.Slug, id)
	}
	if dto.FirstName != "" {
		t.FirstName = dto.FirstName
	}
	if dto.LastName != "" {
		t.LastName = dto.LastName
	}
	t.MiddleName = dto.MiddleName
	t.Bio = dto.Bio
	if err := h.db.Save(&t).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	h.db.Preload("Department").First(&t, t.ID)
	return c.JSON(t)
}

func (h *TeachersHandler) Delete(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid id"})
	}
	var t models.Teacher
	if err := h.db.First(&t, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "not found"})
	}
	if t.PhotoPath != "" {
		_ = h.files.Remove(t.PhotoPath)
	}
	// Публикации преподавателя удаляются каскадом (ON DELETE CASCADE).
	h.db.Delete(&models.Teacher{}, id)
	return c.SendStatus(204)
}

func (h *TeachersHandler) UploadPhoto(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid id"})
	}
	var t models.Teacher
	if err := h.db.First(&t, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "not found"})
	}
	fh, err := c.FormFile("photo")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "photo file required"})
	}
	if !isImageFile(fh) {
		return c.Status(400).JSON(fiber.Map{"error": "only images allowed"})
	}
	ext := strings.ToLower(filepath.Ext(fh.Filename))
	rel := fmt.Sprintf("teachers/%d/photo%s", id, ext)
	if err := h.files.SaveUpload(fh, rel); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	if t.PhotoPath != "" && t.PhotoPath != rel {
		_ = h.files.Remove(t.PhotoPath)
	}
	h.db.Model(&t).Update("photo_path", rel)
	return c.JSON(fiber.Map{"photo_path": rel})
}

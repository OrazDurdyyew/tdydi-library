package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"library/internal/models"
)

type DepartmentsHandler struct {
	db *gorm.DB
}

func NewDepartmentsHandler(db *gorm.DB) *DepartmentsHandler {
	return &DepartmentsHandler{db: db}
}

// reservedDepartmentSlugs — slug'и, конфликтующие с публичными литеральными роутами
// (/publications/teachers, /admin/.../new и т.п.).
var reservedDepartmentSlugs = map[string]bool{
	"teachers": true,
	"new":      true,
	"slug":     true,
}

type DepartmentDTO struct {
	Slug   string `json:"slug"`
	NameTM string `json:"name_tm"`
	NameRU string `json:"name_ru"`
	NameEN string `json:"name_en"`
	DescTM string `json:"description_tm"`
	DescRU string `json:"description_ru"`
	DescEN string `json:"description_en"`
}

// departmentWithCount — кафедра + число преподавателей (для списка в админке).
type departmentWithCount struct {
	models.Department
	TeacherCount int64 `json:"teacher_count"`
}

func (h *DepartmentsHandler) List(c *fiber.Ctx) error {
	var items []departmentWithCount
	h.db.Table("departments").
		Select("departments.*, COUNT(teachers.id) AS teacher_count").
		Joins("LEFT JOIN teachers ON teachers.department_id = departments.id").
		Group("departments.id").
		Order("departments.name_tm ASC").
		Scan(&items)
	return c.JSON(fiber.Map{"data": items})
}

func (h *DepartmentsHandler) Get(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid id"})
	}
	var d models.Department
	if err := h.db.First(&d, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "not found"})
	}
	return c.JSON(d)
}

func (h *DepartmentsHandler) GetBySlug(c *fiber.Ctx) error {
	var d models.Department
	if err := h.db.Where("slug = ?", c.Params("slug")).First(&d).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "not found"})
	}
	return c.JSON(d)
}

func (h *DepartmentsHandler) Create(c *fiber.Ctx) error {
	var dto DepartmentDTO
	if err := c.BodyParser(&dto); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	if dto.NameTM == "" {
		return c.Status(400).JSON(fiber.Map{"error": "name_tm required"})
	}
	if dto.Slug == "" {
		dto.Slug = slugify(dto.NameTM)
	}
	if reservedDepartmentSlugs[dto.Slug] {
		return c.Status(400).JSON(fiber.Map{"error": "slug is reserved"})
	}
	d := models.Department{
		Slug:   dto.Slug,
		NameTM: dto.NameTM,
		NameRU: dto.NameRU,
		NameEN: dto.NameEN,
		DescTM: dto.DescTM,
		DescRU: dto.DescRU,
		DescEN: dto.DescEN,
	}
	if err := h.db.Create(&d).Error; err != nil {
		return c.Status(409).JSON(fiber.Map{"error": "slug already exists"})
	}
	return c.Status(201).JSON(d)
}

func (h *DepartmentsHandler) Update(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid id"})
	}
	var d models.Department
	if err := h.db.First(&d, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "not found"})
	}
	var dto DepartmentDTO
	if err := c.BodyParser(&dto); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	if dto.Slug != "" {
		if reservedDepartmentSlugs[dto.Slug] {
			return c.Status(400).JSON(fiber.Map{"error": "slug is reserved"})
		}
		d.Slug = dto.Slug
	}
	d.NameTM = dto.NameTM
	d.NameRU = dto.NameRU
	d.NameEN = dto.NameEN
	d.DescTM = dto.DescTM
	d.DescRU = dto.DescRU
	d.DescEN = dto.DescEN
	if err := h.db.Save(&d).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(d)
}

func (h *DepartmentsHandler) Delete(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid id"})
	}
	var count int64
	h.db.Model(&models.Teacher{}).Where("department_id = ?", id).Count(&count)
	if count > 0 {
		return c.Status(409).JSON(fiber.Map{"error": "has teachers, delete them first"})
	}
	h.db.Delete(&models.Department{}, id)
	return c.SendStatus(204)
}

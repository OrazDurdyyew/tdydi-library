package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"library/internal/models"
)

type PublicationCategoriesHandler struct {
	db *gorm.DB
}

func NewPublicationCategoriesHandler(db *gorm.DB) *PublicationCategoriesHandler {
	return &PublicationCategoriesHandler{db: db}
}

type PublicationCategoryDTO struct {
	Slug   string `json:"slug"`
	NameTM string `json:"name_tm"`
	NameRU string `json:"name_ru"`
	NameEN string `json:"name_en"`
}

func (h *PublicationCategoriesHandler) List(c *fiber.Ctx) error {
	var items []models.PublicationCategory
	h.db.Order("name_tm ASC").Find(&items)
	return c.JSON(fiber.Map{"data": items})
}

func (h *PublicationCategoriesHandler) Get(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid id"})
	}
	var cat models.PublicationCategory
	if err := h.db.First(&cat, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "not found"})
	}
	return c.JSON(cat)
}

func (h *PublicationCategoriesHandler) Create(c *fiber.Ctx) error {
	var dto PublicationCategoryDTO
	if err := c.BodyParser(&dto); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	if dto.NameTM == "" {
		return c.Status(400).JSON(fiber.Map{"error": "name_tm required"})
	}
	if dto.Slug == "" {
		dto.Slug = slugify(dto.NameTM)
	}
	cat := models.PublicationCategory{
		Slug:   dto.Slug,
		NameTM: dto.NameTM,
		NameRU: dto.NameRU,
		NameEN: dto.NameEN,
	}
	if err := h.db.Create(&cat).Error; err != nil {
		return c.Status(409).JSON(fiber.Map{"error": "slug already exists"})
	}
	return c.Status(201).JSON(cat)
}

func (h *PublicationCategoriesHandler) Update(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid id"})
	}
	var cat models.PublicationCategory
	if err := h.db.First(&cat, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "not found"})
	}
	var dto PublicationCategoryDTO
	if err := c.BodyParser(&dto); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	if dto.Slug != "" {
		cat.Slug = dto.Slug
	}
	cat.NameTM = dto.NameTM
	cat.NameRU = dto.NameRU
	cat.NameEN = dto.NameEN
	if err := h.db.Save(&cat).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(cat)
}

func (h *PublicationCategoriesHandler) Delete(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid id"})
	}
	var count int64
	h.db.Model(&models.Publication{}).Where("category_id = ?", id).Count(&count)
	if count > 0 {
		return c.Status(409).JSON(fiber.Map{"error": "category in use by publications"})
	}
	h.db.Delete(&models.PublicationCategory{}, id)
	return c.SendStatus(204)
}

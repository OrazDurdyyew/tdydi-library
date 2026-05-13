package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"library/internal/models"
)

type FAQHandler struct {
	db *gorm.DB
}

func NewFAQHandler(db *gorm.DB) *FAQHandler {
	return &FAQHandler{db: db}
}

type FAQCategoryDTO struct {
	SortOrder int    `json:"sort_order"`
	NameTM    string `json:"name_tm"`
	NameRU    string `json:"name_ru"`
	NameEN    string `json:"name_en"`
}

type FAQItemDTO struct {
	CategoryID int64  `json:"category_id"`
	SortOrder  int    `json:"sort_order"`
	QuestionTM string `json:"question_tm"`
	QuestionRU string `json:"question_ru"`
	QuestionEN string `json:"question_en"`
	AnswerTM   string `json:"answer_tm"`
	AnswerRU   string `json:"answer_ru"`
	AnswerEN   string `json:"answer_en"`
}

// ---- Public ----

func (h *FAQHandler) ListCategories(c *fiber.Ctx) error {
	var cats []models.FAQCategory
	h.db.Preload("Items", func(db *gorm.DB) *gorm.DB {
		return db.Order("sort_order ASC, id ASC")
	}).Order("sort_order ASC, id ASC").Find(&cats)
	return c.JSON(cats)
}

// ---- Admin: categories ----

func (h *FAQHandler) CreateCategory(c *fiber.Ctx) error {
	var dto FAQCategoryDTO
	if err := c.BodyParser(&dto); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	if dto.NameTM == "" {
		return c.Status(400).JSON(fiber.Map{"error": "name_tm required"})
	}
	cat := models.FAQCategory{
		SortOrder: dto.SortOrder,
		NameTM:    dto.NameTM,
		NameRU:    dto.NameRU,
		NameEN:    dto.NameEN,
	}
	h.db.Create(&cat)
	return c.Status(201).JSON(cat)
}

func (h *FAQHandler) UpdateCategory(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid id"})
	}
	var cat models.FAQCategory
	if err := h.db.First(&cat, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "not found"})
	}
	var dto FAQCategoryDTO
	if err := c.BodyParser(&dto); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	cat.SortOrder = dto.SortOrder
	cat.NameTM = dto.NameTM
	cat.NameRU = dto.NameRU
	cat.NameEN = dto.NameEN
	h.db.Save(&cat)
	return c.JSON(cat)
}

func (h *FAQHandler) DeleteCategory(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid id"})
	}
	var count int64
	h.db.Model(&models.FAQItem{}).Where("category_id = ?", id).Count(&count)
	if count > 0 {
		return c.Status(409).JSON(fiber.Map{"error": "category has items, delete them first"})
	}
	h.db.Delete(&models.FAQCategory{}, id)
	return c.SendStatus(204)
}

// ---- Admin: items ----

func (h *FAQHandler) CreateItem(c *fiber.Ctx) error {
	var dto FAQItemDTO
	if err := c.BodyParser(&dto); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	if dto.QuestionTM == "" || dto.AnswerTM == "" {
		return c.Status(400).JSON(fiber.Map{"error": "question_tm and answer_tm required"})
	}
	var count int64
	h.db.Model(&models.FAQCategory{}).Where("id = ?", dto.CategoryID).Count(&count)
	if count == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "category not found"})
	}
	item := models.FAQItem{
		CategoryID: dto.CategoryID,
		SortOrder:  dto.SortOrder,
		QuestionTM: dto.QuestionTM,
		QuestionRU: dto.QuestionRU,
		QuestionEN: dto.QuestionEN,
		AnswerTM:   dto.AnswerTM,
		AnswerRU:   dto.AnswerRU,
		AnswerEN:   dto.AnswerEN,
	}
	h.db.Create(&item)
	return c.Status(201).JSON(item)
}

func (h *FAQHandler) UpdateItem(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid id"})
	}
	var item models.FAQItem
	if err := h.db.First(&item, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "not found"})
	}
	var dto FAQItemDTO
	if err := c.BodyParser(&dto); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	if dto.CategoryID != 0 {
		item.CategoryID = dto.CategoryID
	}
	item.SortOrder = dto.SortOrder
	item.QuestionTM = dto.QuestionTM
	item.QuestionRU = dto.QuestionRU
	item.QuestionEN = dto.QuestionEN
	item.AnswerTM = dto.AnswerTM
	item.AnswerRU = dto.AnswerRU
	item.AnswerEN = dto.AnswerEN
	h.db.Save(&item)
	return c.JSON(item)
}

func (h *FAQHandler) DeleteItem(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid id"})
	}
	h.db.Delete(&models.FAQItem{}, id)
	return c.SendStatus(204)
}

type ReorderFAQDTO struct {
	IDs []int64 `json:"ids"`
}

func (h *FAQHandler) ReorderCategories(c *fiber.Ctx) error {
	var dto ReorderFAQDTO
	if err := c.BodyParser(&dto); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	for i, id := range dto.IDs {
		h.db.Model(&models.FAQCategory{}).Where("id = ?", id).Update("sort_order", i)
	}
	return c.SendStatus(204)
}

func (h *FAQHandler) ReorderItems(c *fiber.Ctx) error {
	var dto ReorderFAQDTO
	if err := c.BodyParser(&dto); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	for i, id := range dto.IDs {
		h.db.Model(&models.FAQItem{}).Where("id = ?", id).Update("sort_order", i)
	}
	return c.SendStatus(204)
}

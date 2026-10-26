package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"library/internal/models"
)

type PagesHandler struct {
	db *gorm.DB
}

func NewPagesHandler(db *gorm.DB) *PagesHandler {
	return &PagesHandler{db: db}
}

type PageDTO struct {
	Slug      string `json:"slug"`
	InMenu    bool   `json:"in_menu"`
	MenuOrder int    `json:"menu_order"`
	TitleTM   string `json:"title_tm"`
	TitleRU   string `json:"title_ru"`
	TitleEN   string `json:"title_en"`
	ContentTM string `json:"content_tm"`
	ContentRU string `json:"content_ru"`
	ContentEN string `json:"content_en"`
}

// ---- Public ----

func (h *PagesHandler) List(c *fiber.Ctx) error {
	q := h.db.Model(&models.Page{}).Order("menu_order ASC, id ASC")
	if c.Query("in_menu") == "1" {
		q = q.Where("in_menu = ?", true)
	}
	var pages []models.Page
	q.Select("id, slug, in_menu, menu_order, title_tm, title_ru, title_en, created_at, updated_at").Find(&pages)
	return c.JSON(pages)
}

func (h *PagesHandler) Get(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid id"})
	}
	var page models.Page
	if err := h.db.First(&page, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "not found"})
	}
	return c.JSON(page)
}

func (h *PagesHandler) GetBySlug(c *fiber.Ctx) error {
	var page models.Page
	if err := h.db.Where("slug = ?", c.Params("slug")).First(&page).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "not found"})
	}
	return c.JSON(page)
}

// ---- Admin ----

func (h *PagesHandler) Create(c *fiber.Ctx) error {
	var dto PageDTO
	if err := c.BodyParser(&dto); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	if dto.TitleTM == "" {
		return c.Status(400).JSON(fiber.Map{"error": "title_tm required"})
	}
	if dto.Slug == "" {
		dto.Slug = slugify(dto.TitleTM)
	}
	page := models.Page{
		Slug:      dto.Slug,
		InMenu:    dto.InMenu,
		MenuOrder: dto.MenuOrder,
		TitleTM:   dto.TitleTM,
		TitleRU:   dto.TitleRU,
		TitleEN:   dto.TitleEN,
		ContentTM: dto.ContentTM,
		ContentRU: dto.ContentRU,
		ContentEN: dto.ContentEN,
	}
	if err := h.db.Create(&page).Error; err != nil {
		return c.Status(409).JSON(fiber.Map{"error": "slug already exists"})
	}
	return c.Status(201).JSON(page)
}

func (h *PagesHandler) Update(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid id"})
	}
	var page models.Page
	if err := h.db.First(&page, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "not found"})
	}
	var dto PageDTO
	if err := c.BodyParser(&dto); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	if dto.Slug != "" {
		page.Slug = dto.Slug
	}
	page.InMenu = dto.InMenu
	page.MenuOrder = dto.MenuOrder
	page.TitleTM = dto.TitleTM
	page.TitleRU = dto.TitleRU
	page.TitleEN = dto.TitleEN
	page.ContentTM = dto.ContentTM
	page.ContentRU = dto.ContentRU
	page.ContentEN = dto.ContentEN
	h.db.Save(&page)
	return c.JSON(page)
}

func (h *PagesHandler) Delete(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid id"})
	}
	h.db.Delete(&models.Page{}, id)
	return c.SendStatus(204)
}

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

type CollectionsHandler struct {
	db    *gorm.DB
	files *services.FileStorage
}

func NewCollectionsHandler(db *gorm.DB, files *services.FileStorage) *CollectionsHandler {
	return &CollectionsHandler{db: db, files: files}
}

type CollectionDTO struct {
	Slug       string `json:"slug"`
	SortOrder  int    `json:"sort_order"`
	ShowOnHome bool   `json:"show_on_home"`
	NameTM     string `json:"name_tm"`
	NameRU     string `json:"name_ru"`
	NameEN     string `json:"name_en"`
	DescTM     string `json:"description_tm"`
	DescRU     string `json:"description_ru"`
	DescEN     string `json:"description_en"`
}

// ---- Public ----

func (h *CollectionsHandler) List(c *fiber.Ctx) error {
	q := h.db.Model(&models.Collection{}).Order("sort_order ASC, id ASC")
	if c.Query("home") == "1" {
		q = q.Where("show_on_home = ?", true)
	}
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	var total int64
	q.Count(&total)
	var cols []models.Collection
	q.Offset((page - 1) * limit).Limit(limit).Find(&cols)
	return c.JSON(fiber.Map{"data": cols, "total": total})
}

// booksPreload — Preload только опубликованных книг для публичных запросов;
// в админке (admin_id выставлен middleware'ом) возвращаются все.
func (h *CollectionsHandler) booksPreload(c *fiber.Ctx) func(*gorm.DB) *gorm.DB {
	if c.Locals("admin_id") == nil {
		return func(db *gorm.DB) *gorm.DB { return db.Where("books.is_published = ?", true) }
	}
	return func(db *gorm.DB) *gorm.DB { return db }
}

func (h *CollectionsHandler) Get(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid id"})
	}
	var col models.Collection
	if err := h.db.Preload("Books", h.booksPreload(c)).First(&col, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "not found"})
	}
	return c.JSON(col)
}

func (h *CollectionsHandler) GetBySlug(c *fiber.Ctx) error {
	var col models.Collection
	if err := h.db.Preload("Books", h.booksPreload(c)).Where("slug = ?", c.Params("slug")).First(&col).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "not found"})
	}
	return c.JSON(col)
}

// ---- Admin ----

func (h *CollectionsHandler) Create(c *fiber.Ctx) error {
	var dto CollectionDTO
	if err := c.BodyParser(&dto); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	if dto.NameTM == "" {
		return c.Status(400).JSON(fiber.Map{"error": "name_tm required"})
	}
	if dto.Slug == "" {
		dto.Slug = slugify(dto.NameTM)
	}
	col := models.Collection{
		Slug:       dto.Slug,
		SortOrder:  dto.SortOrder,
		ShowOnHome: dto.ShowOnHome,
		NameTM:     dto.NameTM,
		NameRU:     dto.NameRU,
		NameEN:     dto.NameEN,
		DescTM:     dto.DescTM,
		DescRU:     dto.DescRU,
		DescEN:     dto.DescEN,
	}
	if err := h.db.Create(&col).Error; err != nil {
		return c.Status(409).JSON(fiber.Map{"error": "slug already exists"})
	}
	return c.Status(201).JSON(col)
}

func (h *CollectionsHandler) Update(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid id"})
	}
	var col models.Collection
	if err := h.db.First(&col, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "not found"})
	}
	var dto CollectionDTO
	if err := c.BodyParser(&dto); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	if dto.Slug != "" {
		col.Slug = dto.Slug
	}
	col.SortOrder = dto.SortOrder
	col.ShowOnHome = dto.ShowOnHome
	col.NameTM = dto.NameTM
	col.NameRU = dto.NameRU
	col.NameEN = dto.NameEN
	col.DescTM = dto.DescTM
	col.DescRU = dto.DescRU
	col.DescEN = dto.DescEN
	h.db.Save(&col)
	return c.JSON(col)
}

func (h *CollectionsHandler) Delete(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid id"})
	}
	h.db.Where("collection_id = ?", id).Delete(&models.CollectionBook{})
	h.db.Delete(&models.Collection{}, id)
	return c.SendStatus(204)
}

func (h *CollectionsHandler) UploadCover(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid id"})
	}
	var col models.Collection
	if err := h.db.First(&col, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "not found"})
	}
	fh, err := c.FormFile("cover")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "cover file required"})
	}
	if !isImageFile(fh) {
		return c.Status(400).JSON(fiber.Map{"error": "only images allowed"})
	}
	ext := strings.ToLower(filepath.Ext(fh.Filename))
	rel := fmt.Sprintf("collections/%d/cover%s", id, ext)
	if err := h.files.SaveUpload(fh, rel); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	h.db.Model(&col).Update("cover_path", rel)
	return c.JSON(fiber.Map{"cover_path": rel})
}

type AddBookDTO struct {
	BookID    int64 `json:"book_id"`
	SortOrder int   `json:"sort_order"`
}

func (h *CollectionsHandler) AddBook(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid id"})
	}
	var dto AddBookDTO
	if err := c.BodyParser(&dto); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	cb := models.CollectionBook{CollectionID: id, BookID: dto.BookID, SortOrder: dto.SortOrder}
	if err := h.db.Create(&cb).Error; err != nil {
		// already exists — update sort_order
		h.db.Model(&cb).Where("collection_id = ? AND book_id = ?", id, dto.BookID).
			Update("sort_order", dto.SortOrder)
	}
	return c.SendStatus(204)
}

func (h *CollectionsHandler) RemoveBook(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid collection id"})
	}
	bookID, err := strconv.ParseInt(c.Params("book_id"), 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid book_id"})
	}
	h.db.Where("collection_id = ? AND book_id = ?", id, bookID).Delete(&models.CollectionBook{})
	return c.SendStatus(204)
}

type ReorderBooksDTO struct {
	BookIDs []int64 `json:"book_ids"`
}

func (h *CollectionsHandler) ReorderBooks(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid id"})
	}
	var dto ReorderBooksDTO
	if err := c.BodyParser(&dto); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	for i, bookID := range dto.BookIDs {
		h.db.Model(&models.CollectionBook{}).
			Where("collection_id = ? AND book_id = ?", id, bookID).
			Update("sort_order", i)
	}
	return c.SendStatus(204)
}

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

type HomeHandler struct {
	db    *gorm.DB
	files *services.FileStorage
}

func NewHomeHandler(db *gorm.DB, files *services.FileStorage) *HomeHandler {
	return &HomeHandler{db: db, files: files}
}

// GetHomeData — public endpoint: all home-page building blocks in one response.
func (h *HomeHandler) GetHomeData(c *fiber.Ctx) error {
	var slides []models.HomeSlide
	h.db.Where("is_active = ?", true).Order("sort_order ASC").Find(&slides)

	var homeSections []models.HomeSection
	h.db.Preload("Section").Where("is_active = ?", true).Order("sort_order ASC").Find(&homeSections)

	var homeCollections []models.HomeCollection
	h.db.Preload("Collection").Where("is_active = ?", true).Order("sort_order ASC").Find(&homeCollections)

	return c.JSON(fiber.Map{
		"slides":      slides,
		"sections":    homeSections,
		"collections": homeCollections,
	})
}

// ---- Slides admin ----

type SlideDTO struct {
	SortOrder  int    `json:"sort_order"`
	IsActive   bool   `json:"is_active"`
	LinkURL    string `json:"link_url"`
	TitleTM    string `json:"title_tm"`
	TitleRU    string `json:"title_ru"`
	TitleEN    string `json:"title_en"`
	SubtitleTM string `json:"subtitle_tm"`
	SubtitleRU string `json:"subtitle_ru"`
	SubtitleEN string `json:"subtitle_en"`
}

func (h *HomeHandler) ListSlides(c *fiber.Ctx) error {
	var slides []models.HomeSlide
	h.db.Order("sort_order ASC, id ASC").Find(&slides)
	return c.JSON(slides)
}

func (h *HomeHandler) CreateSlide(c *fiber.Ctx) error {
	var dto SlideDTO
	if err := c.BodyParser(&dto); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	slide := models.HomeSlide{
		SortOrder:  dto.SortOrder,
		IsActive:   dto.IsActive,
		LinkURL:    dto.LinkURL,
		TitleTM:    dto.TitleTM,
		TitleRU:    dto.TitleRU,
		TitleEN:    dto.TitleEN,
		SubtitleTM: dto.SubtitleTM,
		SubtitleRU: dto.SubtitleRU,
		SubtitleEN: dto.SubtitleEN,
	}
	h.db.Create(&slide)
	return c.Status(201).JSON(slide)
}

func (h *HomeHandler) UpdateSlide(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid id"})
	}
	var slide models.HomeSlide
	if err := h.db.First(&slide, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "not found"})
	}
	var dto SlideDTO
	if err := c.BodyParser(&dto); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	slide.SortOrder = dto.SortOrder
	slide.IsActive = dto.IsActive
	slide.LinkURL = dto.LinkURL
	slide.TitleTM = dto.TitleTM
	slide.TitleRU = dto.TitleRU
	slide.TitleEN = dto.TitleEN
	slide.SubtitleTM = dto.SubtitleTM
	slide.SubtitleRU = dto.SubtitleRU
	slide.SubtitleEN = dto.SubtitleEN
	h.db.Save(&slide)
	return c.JSON(slide)
}

func (h *HomeHandler) DeleteSlide(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid id"})
	}
	h.db.Delete(&models.HomeSlide{}, id)
	return c.SendStatus(204)
}

func (h *HomeHandler) UploadSlideImage(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid id"})
	}
	var slide models.HomeSlide
	if err := h.db.First(&slide, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "not found"})
	}
	fh, err := c.FormFile("image")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "image file required"})
	}
	if !isImageFile(fh) {
		return c.Status(400).JSON(fiber.Map{"error": "only images allowed"})
	}
	ext := strings.ToLower(filepath.Ext(fh.Filename))
	rel := fmt.Sprintf("slides/%d/image%s", id, ext)
	if err := h.files.SaveUpload(fh, rel); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	h.db.Model(&slide).Update("image_path", rel)
	return c.JSON(fiber.Map{"image_path": rel})
}

type ReorderSlidesDTO struct {
	IDs []int64 `json:"ids"`
}

func (h *HomeHandler) ReorderSlides(c *fiber.Ctx) error {
	var dto ReorderSlidesDTO
	if err := c.BodyParser(&dto); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	for i, id := range dto.IDs {
		h.db.Model(&models.HomeSlide{}).Where("id = ?", id).Update("sort_order", i)
	}
	return c.SendStatus(204)
}

// ---- Home sections admin ----

type HomeSectionDTO struct {
	SectionID int64 `json:"section_id"`
	SortOrder int   `json:"sort_order"`
	IsActive  bool  `json:"is_active"`
}

func (h *HomeHandler) ListHomeSections(c *fiber.Ctx) error {
	var items []models.HomeSection
	h.db.Preload("Section").Order("sort_order ASC, id ASC").Find(&items)
	return c.JSON(items)
}

func (h *HomeHandler) AddHomeSection(c *fiber.Ctx) error {
	var dto HomeSectionDTO
	if err := c.BodyParser(&dto); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	item := models.HomeSection{SectionID: dto.SectionID, SortOrder: dto.SortOrder, IsActive: dto.IsActive}
	if err := h.db.Create(&item).Error; err != nil {
		return c.Status(409).JSON(fiber.Map{"error": "section already on home"})
	}
	h.db.Preload("Section").First(&item, item.ID)
	return c.Status(201).JSON(item)
}

func (h *HomeHandler) UpdateHomeSection(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid id"})
	}
	var item models.HomeSection
	if err := h.db.First(&item, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "not found"})
	}
	var dto HomeSectionDTO
	if err := c.BodyParser(&dto); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	item.SortOrder = dto.SortOrder
	item.IsActive = dto.IsActive
	h.db.Save(&item)
	return c.JSON(item)
}

func (h *HomeHandler) RemoveHomeSection(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid id"})
	}
	h.db.Delete(&models.HomeSection{}, id)
	return c.SendStatus(204)
}

// ---- Home collections admin ----

type HomeCollectionDTO struct {
	CollectionID int64 `json:"collection_id"`
	SortOrder    int   `json:"sort_order"`
	IsActive     bool  `json:"is_active"`
}

func (h *HomeHandler) ListHomeCollections(c *fiber.Ctx) error {
	var items []models.HomeCollection
	h.db.Preload("Collection").Order("sort_order ASC, id ASC").Find(&items)
	return c.JSON(items)
}

func (h *HomeHandler) AddHomeCollection(c *fiber.Ctx) error {
	var dto HomeCollectionDTO
	if err := c.BodyParser(&dto); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	item := models.HomeCollection{CollectionID: dto.CollectionID, SortOrder: dto.SortOrder, IsActive: dto.IsActive}
	if err := h.db.Create(&item).Error; err != nil {
		return c.Status(409).JSON(fiber.Map{"error": "collection already on home"})
	}
	h.db.Preload("Collection").First(&item, item.ID)
	return c.Status(201).JSON(item)
}

func (h *HomeHandler) UpdateHomeCollection(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid id"})
	}
	var item models.HomeCollection
	if err := h.db.First(&item, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "not found"})
	}
	var dto HomeCollectionDTO
	if err := c.BodyParser(&dto); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	item.SortOrder = dto.SortOrder
	item.IsActive = dto.IsActive
	h.db.Save(&item)
	return c.JSON(item)
}

func (h *HomeHandler) RemoveHomeCollection(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid id"})
	}
	h.db.Delete(&models.HomeCollection{}, id)
	return c.SendStatus(204)
}

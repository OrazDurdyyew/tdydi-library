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

// HomeSectionBlock — блок раздела на главной: сам раздел + несколько его книг.
type HomeSectionBlock struct {
	Section models.Section `json:"section"`
	Books   []models.Book  `json:"books"`
}

// GetHomeData — public endpoint: all home-page building blocks in one response.
// Включает слайды, новинки, блоки разделов с книгами, подборки и последние новости.
func (h *HomeHandler) GetHomeData(c *fiber.Ctx) error {
	var slides []models.HomeSlide
	h.db.Where("is_active = ?", true).Order("sort_order ASC").Find(&slides)

	var newBooks []models.Book
	h.db.Where("is_published = ?", true).Preload("Sections").
		Order("created_at DESC").Limit(12).Find(&newBooks)

	var homeSections []models.HomeSection
	h.db.Preload("Section").Where("is_active = ?", true).Order("sort_order ASC").Find(&homeSections)

	sectionBlocks := make([]HomeSectionBlock, 0, len(homeSections))
	for _, hs := range homeSections {
		if hs.Section == nil {
			continue
		}
		var books []models.Book
		h.db.Joins("JOIN book_sections bs ON bs.book_id = books.id").
			Where("bs.section_id = ? AND books.is_published = ?", hs.SectionID, true).
			Order("books.created_at DESC").Limit(6).Find(&books)
		sectionBlocks = append(sectionBlocks, HomeSectionBlock{Section: *hs.Section, Books: books})
	}

	var homeCollections []models.HomeCollection
	h.db.Preload("Collection").Where("is_active = ?", true).Order("sort_order ASC").Find(&homeCollections)

	var latestPosts []models.Post
	h.db.Select("id, slug, cover_path, author, published_at, status, title_tm, title_ru, title_en").
		Where("status = ?", models.PostStatusPublished).
		Order("published_at DESC, created_at DESC").Limit(4).Find(&latestPosts)

	var promos []models.HomePromo
	h.db.Where("is_active = ?", true).Order("sort_order ASC, id ASC").Find(&promos)

	var totalPublications, totalTeachers, totalDepartments int64
	h.db.Model(&models.Publication{}).Where("is_published = ?", true).Count(&totalPublications)
	h.db.Model(&models.Teacher{}).Count(&totalTeachers)
	h.db.Model(&models.Department{}).Count(&totalDepartments)

	return c.JSON(fiber.Map{
		"slides":         slides,
		"new_books":      newBooks,
		"section_blocks": sectionBlocks,
		"collections":    homeCollections,
		"latest_posts":   latestPosts,
		"promos":         promos,
		"publications_stats": fiber.Map{
			"total_publications": totalPublications,
			"total_teachers":     totalTeachers,
			"total_departments":  totalDepartments,
		},
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

// ---- Home promos admin ----

type PromoDTO struct {
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

func (h *HomeHandler) ListPromos(c *fiber.Ctx) error {
	var promos []models.HomePromo
	h.db.Order("sort_order ASC, id ASC").Find(&promos)
	return c.JSON(promos)
}

func (h *HomeHandler) CreatePromo(c *fiber.Ctx) error {
	var dto PromoDTO
	if err := c.BodyParser(&dto); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	promo := models.HomePromo{
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
	h.db.Create(&promo)
	return c.Status(201).JSON(promo)
}

func (h *HomeHandler) UpdatePromo(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid id"})
	}
	var promo models.HomePromo
	if err := h.db.First(&promo, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "not found"})
	}
	var dto PromoDTO
	if err := c.BodyParser(&dto); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	promo.SortOrder = dto.SortOrder
	promo.IsActive = dto.IsActive
	promo.LinkURL = dto.LinkURL
	promo.TitleTM = dto.TitleTM
	promo.TitleRU = dto.TitleRU
	promo.TitleEN = dto.TitleEN
	promo.SubtitleTM = dto.SubtitleTM
	promo.SubtitleRU = dto.SubtitleRU
	promo.SubtitleEN = dto.SubtitleEN
	h.db.Save(&promo)
	return c.JSON(promo)
}

func (h *HomeHandler) DeletePromo(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid id"})
	}
	h.db.Delete(&models.HomePromo{}, id)
	return c.SendStatus(204)
}

func (h *HomeHandler) ReorderPromos(c *fiber.Ctx) error {
	var dto ReorderSlidesDTO
	if err := c.BodyParser(&dto); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	for i, id := range dto.IDs {
		h.db.Model(&models.HomePromo{}).Where("id = ?", id).Update("sort_order", i)
	}
	return c.SendStatus(204)
}

func (h *HomeHandler) UploadPromoImage(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid id"})
	}
	var promo models.HomePromo
	if err := h.db.First(&promo, id).Error; err != nil {
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
	rel := fmt.Sprintf("promos/%d/image%s", id, ext)
	if err := h.files.SaveUpload(fh, rel); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	h.db.Model(&promo).Update("image_path", rel)
	return c.JSON(fiber.Map{"image_path": rel})
}

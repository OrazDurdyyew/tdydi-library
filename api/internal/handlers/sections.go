package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"library/internal/models"
)

type SectionsHandler struct {
	db *gorm.DB
}

func NewSectionsHandler(db *gorm.DB) *SectionsHandler {
	return &SectionsHandler{db: db}
}

type SectionDTO struct {
	Slug      string `json:"slug"`
	ParentID  *int64 `json:"parent_id"`
	SortOrder int    `json:"sort_order"`
	ImagePath string `json:"image_path"`
	NameTM    string `json:"name_tm"`
	NameRU    string `json:"name_ru"`
	NameEN    string `json:"name_en"`
	DescTM    string `json:"description_tm"`
	DescRU    string `json:"description_ru"`
	DescEN    string `json:"description_en"`
}

// List возвращает разделы. Опциональный параметр: ?parent_id=N или ?parent_id=null.
// Для каждого раздела заполняется book_count (книги, привязанные напрямую).
func (h *SectionsHandler) List(c *fiber.Ctx) error {
	q := h.db.Model(&models.Section{}).Order("sort_order ASC, id ASC")
	if p := c.Query("parent_id"); p != "" {
		if p == "null" {
			q = q.Where("parent_id IS NULL")
		} else {
			pid, err := strconv.ParseInt(p, 10, 64)
			if err != nil {
				return c.Status(400).JSON(fiber.Map{"error": "invalid parent_id"})
			}
			q = q.Where("parent_id = ?", pid)
		}
	}
	var sections []models.Section
	if err := q.Find(&sections).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	for i := range sections {
		sections[i].BookCount = h.countBooks(sections[i].ID)
	}
	return c.JSON(fiber.Map{"data": sections})
}

// countBooks возвращает количество книг, привязанных напрямую к разделу.
func (h *SectionsHandler) countBooks(sectionID int64) int64 {
	var n int64
	h.db.Model(&models.BookSection{}).Where("section_id = ?", sectionID).Count(&n)
	return n
}

func (h *SectionsHandler) Get(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid id"})
	}
	var s models.Section
	if err := h.db.Preload("Children", orderBySortOrder).First(&s, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "not found"})
	}
	return c.JSON(s)
}

// GetBySlug — публичная страница раздела: сам раздел (+ родитель для хлебных крошек),
// его подразделы (карточки с book_count) и книги, привязанные непосредственно к разделу.
func (h *SectionsHandler) GetBySlug(c *fiber.Ctx) error {
	slug := c.Params("slug")
	var s models.Section
	if err := h.db.Preload("Children", orderBySortOrder).Preload("Parent").
		Where("slug = ?", slug).First(&s).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "not found"})
	}

	s.BookCount = h.countBooks(s.ID)
	for i := range s.Children {
		s.Children[i].BookCount = h.countBooks(s.Children[i].ID)
	}

	bq := h.db.Model(&models.Book{}).
		Joins("JOIN book_sections bs ON bs.book_id = books.id").
		Where("bs.section_id = ?", s.ID)
	if c.Locals("admin_id") == nil {
		bq = bq.Where("books.is_published = ?", true)
	}
	var books []models.Book
	bq.Preload("Sections").Order("books.created_at DESC").Find(&books)

	return c.JSON(fiber.Map{"section": s, "books": books})
}

// Reorder — drag-and-drop переупорядочивание разделов в админке.
// Принимает {items: [{id, parent_id, sort_order}, ...]}.
func (h *SectionsHandler) Reorder(c *fiber.Ctx) error {
	var body struct {
		Items []struct {
			ID        int64  `json:"id"`
			ParentID  *int64 `json:"parent_id"`
			SortOrder int    `json:"sort_order"`
		} `json:"items"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
	}
	err := h.db.Transaction(func(tx *gorm.DB) error {
		for _, it := range body.Items {
			if err := tx.Model(&models.Section{}).Where("id = ?", it.ID).
				Updates(map[string]interface{}{
					"parent_id":  it.ParentID,
					"sort_order": it.SortOrder,
				}).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"status": "ok"})
}

func (h *SectionsHandler) Create(c *fiber.Ctx) error {
	dto, err := parseSectionDTO(c)
	if err != nil {
		return err
	}
	s := dto.toModel()
	if err := h.db.Create(&s).Error; err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(201).JSON(s)
}

func (h *SectionsHandler) Update(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid id"})
	}
	dto, err := parseSectionDTO(c)
	if err != nil {
		return err
	}
	updates := map[string]interface{}{
		"slug":           dto.Slug,
		"parent_id":      dto.ParentID,
		"sort_order":     dto.SortOrder,
		"image_path":     dto.ImagePath,
		"name_tm":        dto.NameTM,
		"name_ru":        dto.NameRU,
		"name_en":        dto.NameEN,
		"description_tm": dto.DescTM,
		"description_ru": dto.DescRU,
		"description_en": dto.DescEN,
	}
	if err := h.db.Model(&models.Section{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	var s models.Section
	if err := h.db.First(&s, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "not found"})
	}
	return c.JSON(s)
}

func (h *SectionsHandler) Delete(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid id"})
	}

	var bookCount int64
	h.db.Model(&models.BookSection{}).Where("section_id = ?", id).Count(&bookCount)
	if bookCount > 0 {
		return c.Status(409).JSON(fiber.Map{
			"error":       "section has books, remove them first",
			"books_count": bookCount,
		})
	}
	var childCount int64
	h.db.Model(&models.Section{}).Where("parent_id = ?", id).Count(&childCount)
	if childCount > 0 {
		return c.Status(409).JSON(fiber.Map{
			"error":          "section has children, remove them first",
			"children_count": childCount,
		})
	}

	if err := h.db.Delete(&models.Section{}, id).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"status": "deleted"})
}

func parseSectionDTO(c *fiber.Ctx) (SectionDTO, error) {
	var dto SectionDTO
	if err := c.BodyParser(&dto); err != nil {
		return dto, c.Status(400).JSON(fiber.Map{"error": "invalid request"})
	}
	if dto.Slug == "" || dto.NameTM == "" || dto.NameRU == "" || dto.NameEN == "" {
		return dto, c.Status(400).JSON(fiber.Map{"error": "slug, name_tm, name_ru, name_en are required"})
	}
	return dto, nil
}

func (d SectionDTO) toModel() models.Section {
	return models.Section{
		Slug:      d.Slug,
		ParentID:  d.ParentID,
		SortOrder: d.SortOrder,
		ImagePath: d.ImagePath,
		NameTM:    d.NameTM,
		NameRU:    d.NameRU,
		NameEN:    d.NameEN,
		DescTM:    d.DescTM,
		DescRU:    d.DescRU,
		DescEN:    d.DescEN,
	}
}

func orderBySortOrder(db *gorm.DB) *gorm.DB {
	return db.Order("sort_order ASC, id ASC")
}

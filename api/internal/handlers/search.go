package handlers

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"library/internal/models"
)

type SearchHandler struct {
	db *gorm.DB
}

func NewSearchHandler(db *gorm.DB) *SearchHandler {
	return &SearchHandler{db: db}
}

// Get — GET /api/search?q=... Полнотекстовый поиск (PostgreSQL FTS) по книгам,
// периодике, новостям, FAQ и статичным страницам. Результаты разделены по типам.
func (h *SearchHandler) Get(c *fiber.Ctx) error {
	q := strings.TrimSpace(c.Query("q"))

	empty := fiber.Map{
		"books":       []any{},
		"periodicals": []any{},
		"posts":       []any{},
		"faq":         []any{},
		"pages":       []any{},
	}
	if len([]rune(q)) < 2 {
		return c.JSON(empty)
	}

	const match = "search_vector @@ websearch_to_tsquery('simple', unaccent(?))"

	var books []models.Book
	h.db.Where("is_published = ? AND "+match, true, q).
		Preload("Sections").Limit(12).Find(&books)

	var periodicals []models.Periodical
	h.db.Where(match, q).Limit(12).Find(&periodicals)

	var posts []models.Post
	h.db.Select("id, slug, cover_path, author, published_at, status, title_tm, title_ru, title_en").
		Where("status = ? AND "+match, models.PostStatusPublished, q).
		Limit(12).Find(&posts)

	var faq []models.FAQItem
	h.db.Where(match, q).Limit(12).Find(&faq)

	var pages []models.Page
	h.db.Select("id, slug, in_menu, title_tm, title_ru, title_en").
		Where(match, q).Limit(12).Find(&pages)

	return c.JSON(fiber.Map{
		"books":       books,
		"periodicals": periodicals,
		"posts":       posts,
		"faq":         faq,
		"pages":       pages,
	})
}

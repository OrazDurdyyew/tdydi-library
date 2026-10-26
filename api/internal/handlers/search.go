package handlers

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

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
//
// Гибридная стратегия:
//   - FTS (websearch_to_tsquery) ловит полные слова и фразы;
//   - ILIKE по unaccent(lower(...)) ловит подстроки/префиксы ("tes" → "test"),
//     регистронезависимо, без учёта диакритики.
//
// При limit=N (опционально) возвращается срез N книг — используется в live-подсказках
// в шапке (без поиска по остальным сущностям).
func (h *SearchHandler) Get(c *fiber.Ctx) error {
	q := strings.TrimSpace(c.Query("q"))
	booksOnly := c.Query("books_only") == "1"
	limit := 12
	if booksOnly {
		limit = 8
	}

	empty := fiber.Map{
		"books":        []any{},
		"periodicals":  []any{},
		"posts":        []any{},
		"faq":          []any{},
		"pages":        []any{},
		"publications": []any{},
	}
	if len([]rune(q)) < 2 {
		return c.JSON(empty)
	}

	like := "%" + strings.ToLower(q) + "%"

	bookWhere := "(search_vector @@ websearch_to_tsquery('simple', unaccent(?))" +
		" OR unaccent(lower(title)) LIKE unaccent(?)" +
		" OR unaccent(lower(authors)) LIKE unaccent(?)" +
		" OR unaccent(lower(annotation)) LIKE unaccent(?)" +
		" OR lower(isbn) LIKE ?" +
		" OR lower(publisher) LIKE ?)"

	rankOrder := clause.Expr{
		SQL:  "ts_rank(search_vector, websearch_to_tsquery('simple', unaccent(?))) DESC, created_at DESC",
		Vars: []any{q},
	}

	var books []models.Book
	h.db.Where("is_published = ? AND "+bookWhere, true, q, like, like, like, like, like).
		Preload("Sections").
		Order(rankOrder).
		Limit(limit).Find(&books)

	if booksOnly {
		return c.JSON(fiber.Map{
			"books":        books,
			"periodicals":  []any{},
			"posts":        []any{},
			"faq":          []any{},
			"pages":        []any{},
			"publications": []any{},
		})
	}

	periodicalWhere := "(search_vector @@ websearch_to_tsquery('simple', unaccent(?))" +
		" OR unaccent(lower(name_tm)) LIKE unaccent(?)" +
		" OR unaccent(lower(name_ru)) LIKE unaccent(?)" +
		" OR unaccent(lower(name_en)) LIKE unaccent(?))"
	var periodicals []models.Periodical
	h.db.Where(periodicalWhere, q, like, like, like).Limit(12).Find(&periodicals)

	postWhere := "(search_vector @@ websearch_to_tsquery('simple', unaccent(?))" +
		" OR unaccent(lower(title_tm)) LIKE unaccent(?)" +
		" OR unaccent(lower(title_ru)) LIKE unaccent(?)" +
		" OR unaccent(lower(title_en)) LIKE unaccent(?))"
	var posts []models.Post
	h.db.Select("id, slug, cover_path, author, published_at, status, title_tm, title_ru, title_en").
		Where("status = ? AND "+postWhere, models.PostStatusPublished, q, like, like, like).
		Limit(12).Find(&posts)

	faqWhere := "(search_vector @@ websearch_to_tsquery('simple', unaccent(?))" +
		" OR unaccent(lower(question_tm)) LIKE unaccent(?)" +
		" OR unaccent(lower(question_ru)) LIKE unaccent(?)" +
		" OR unaccent(lower(question_en)) LIKE unaccent(?))"
	var faq []models.FAQItem
	h.db.Where(faqWhere, q, like, like, like).Limit(12).Find(&faq)

	pageWhere := "(search_vector @@ websearch_to_tsquery('simple', unaccent(?))" +
		" OR unaccent(lower(title_tm)) LIKE unaccent(?)" +
		" OR unaccent(lower(title_ru)) LIKE unaccent(?)" +
		" OR unaccent(lower(title_en)) LIKE unaccent(?))"
	var pages []models.Page
	h.db.Select("id, slug, in_menu, title_tm, title_ru, title_en").
		Where(pageWhere, q, like, like, like).Limit(12).Find(&pages)

	publicationWhere := "(search_vector @@ websearch_to_tsquery('simple', unaccent(?))" +
		" OR unaccent(lower(name_tm)) LIKE unaccent(?)" +
		" OR unaccent(lower(name_ru)) LIKE unaccent(?)" +
		" OR unaccent(lower(name_en)) LIKE unaccent(?))"
	var publications []models.Publication
	h.db.Where("is_published = ? AND "+publicationWhere, true, q, like, like, like).
		Preload("Teacher").Preload("Category").
		Limit(12).Find(&publications)

	return c.JSON(fiber.Map{
		"books":        books,
		"periodicals":  periodicals,
		"posts":        posts,
		"faq":          faq,
		"pages":        pages,
		"publications": publications,
	})
}


package handlers

import (
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"library/internal/models"
	"library/internal/services"
)

type StatsHandler struct {
	db    *gorm.DB
	perms *services.PermissionsService
}

func NewStatsHandler(db *gorm.DB, perms *services.PermissionsService) *StatsHandler {
	return &StatsHandler{db: db, perms: perms}
}

type BookStat struct {
	ID             int64  `json:"id"`
	Title          string `json:"title"`
	Authors        string `json:"authors"`
	ViewsCount     int64  `json:"views_count"`
	DownloadsCount int64  `json:"downloads_count"`
}

type DailyActivity struct {
	Date  string `json:"date"`
	Views int64  `json:"views"`
}

type SectionStat struct {
	SectionID   int64  `json:"section_id"`
	SectionName string `json:"section_name"`
	BookCount   int64  `json:"book_count"`
	ViewsTotal  int64  `json:"views_total"`
}

func (h *StatsHandler) Overview(c *fiber.Ctx) error {
	adminID, _ := c.Locals("admin_id").(int64)
	role, _ := c.Locals("admin_role").(string)
	isSuperadmin := role == string(models.RoleSuperadmin)

	days, _ := strconv.Atoi(c.Query("days", "30"))
	if days < 1 || days > 365 {
		days = 30
	}
	since := time.Now().AddDate(0, 0, -days)

	// Determine which book IDs this admin can see
	var visibleBookIDs []int64
	if !isSuperadmin {
		allowedSections, err := h.perms.AllowedSectionIDs(adminID)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		if len(allowedSections) == 0 {
			return c.JSON(emptyStats(days))
		}
		h.db.Table("books").
			Joins("JOIN book_sections bs ON bs.book_id = books.id AND bs.section_id IN ?", allowedSections).
			Pluck("books.id", &visibleBookIDs)
		if len(visibleBookIDs) == 0 {
			return c.JSON(emptyStats(days))
		}
	}

	bookQ := h.db.Model(&models.Book{})
	if !isSuperadmin {
		bookQ = bookQ.Where("id IN ?", visibleBookIDs)
	}

	var totalBooks int64
	bookQ.Count(&totalBooks)

	type Totals struct {
		Views     int64
		Downloads int64
	}
	var totals Totals
	totalsQ := h.db.Model(&models.Book{}).
		Select("COALESCE(SUM(views_count),0) as views, COALESCE(SUM(downloads_count),0) as downloads")
	if !isSuperadmin {
		totalsQ = totalsQ.Where("id IN ?", visibleBookIDs)
	}
	totalsQ.Scan(&totals)

	// Top 10 books by views
	topQ := h.db.Model(&models.Book{}).
		Select("id, title, authors, views_count, downloads_count").
		Order("views_count DESC").Limit(10)
	if !isSuperadmin {
		topQ = topQ.Where("id IN ?", visibleBookIDs)
	}
	var topBooks []BookStat
	topQ.Scan(&topBooks)
	if topBooks == nil {
		topBooks = []BookStat{}
	}

	// Daily view activity
	type RawDay struct {
		Day   string
		Count int64
	}
	viewQ := h.db.Model(&models.ViewLog{}).
		Select("TO_CHAR(created_at, 'YYYY-MM-DD') as day, COUNT(*) as count").
		Where("entity_type = ? AND created_at >= ?", models.ViewEntityBook, since).
		Group("TO_CHAR(created_at, 'YYYY-MM-DD')").
		Order("day ASC")
	if !isSuperadmin {
		viewQ = viewQ.Where("entity_id IN ?", visibleBookIDs)
	}
	var rawDays []RawDay
	viewQ.Scan(&rawDays)

	daily := make([]DailyActivity, len(rawDays))
	for i, d := range rawDays {
		daily[i] = DailyActivity{Date: d.Day, Views: d.Count}
	}

	// Top sections by views (superadmin only)
	var topSections []SectionStat
	if isSuperadmin {
		h.db.Table("sections s").
			Select("s.id as section_id, s.name_tm as section_name, COUNT(DISTINCT bs.book_id) as book_count, COALESCE(SUM(b.views_count),0) as views_total").
			Joins("LEFT JOIN book_sections bs ON bs.section_id = s.id").
			Joins("LEFT JOIN books b ON b.id = bs.book_id").
			Group("s.id, s.name_tm").
			Order("views_total DESC").
			Limit(10).
			Scan(&topSections)
	}
	if topSections == nil {
		topSections = []SectionStat{}
	}

	return c.JSON(fiber.Map{
		"total_books":     totalBooks,
		"total_views":     totals.Views,
		"total_downloads": totals.Downloads,
		"top_books":       topBooks,
		"daily_activity":  daily,
		"top_sections":    topSections,
		"period_days":     days,
	})
}

// ExportCSV exports book stats as CSV (superadmin only).
func (h *StatsHandler) ExportCSV(c *fiber.Ctx) error {
	var books []models.Book
	h.db.Select("id, title, authors, language, year, views_count, downloads_count").
		Order("views_count DESC").Find(&books)

	c.Set("Content-Type", "text/csv; charset=utf-8")
	c.Set("Content-Disposition", "attachment; filename=library-stats.csv")

	out := "ID,Title,Authors,Language,Year,Views,Downloads\n"
	for _, b := range books {
		out += strconv.FormatInt(b.ID, 10) + "," +
			csvEscape(b.Title) + "," +
			csvEscape(b.Authors) + "," +
			b.Language + "," +
			strconv.Itoa(int(b.Year)) + "," +
			strconv.FormatInt(b.ViewsCount, 10) + "," +
			strconv.FormatInt(b.DownloadsCount, 10) + "\n"
	}
	return c.SendString(out)
}

func emptyStats(days int) fiber.Map {
	return fiber.Map{
		"total_books":     0,
		"total_views":     0,
		"total_downloads": 0,
		"top_books":       []BookStat{},
		"daily_activity":  []DailyActivity{},
		"top_sections":    []SectionStat{},
		"period_days":     days,
	}
}

func csvEscape(s string) string {
	needsQuote := false
	for _, ch := range s {
		if ch == '"' || ch == ',' || ch == '\n' || ch == '\r' {
			needsQuote = true
			break
		}
	}
	if !needsQuote {
		return s
	}
	var b strings.Builder
	b.WriteByte('"')
	for _, ch := range s {
		if ch == '"' {
			b.WriteString(`""`)
		} else {
			b.WriteRune(ch)
		}
	}
	b.WriteByte('"')
	return b.String()
}

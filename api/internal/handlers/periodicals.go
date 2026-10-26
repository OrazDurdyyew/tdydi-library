package handlers

import (
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"library/internal/models"
	"library/internal/services"
)

type PeriodicalsHandler struct {
	db    *gorm.DB
	files *services.FileStorage
}

func NewPeriodicalsHandler(db *gorm.DB, files *services.FileStorage) *PeriodicalsHandler {
	return &PeriodicalsHandler{db: db, files: files}
}

type PeriodicalDTO struct {
	Slug        string                `json:"slug"`
	Type        models.PeriodicalType `json:"type"`
	Periodicity string                `json:"periodicity"`
	NameTM      string                `json:"name_tm"`
	NameRU      string                `json:"name_ru"`
	NameEN      string                `json:"name_en"`
	DescTM      string                `json:"description_tm"`
	DescRU      string                `json:"description_ru"`
	DescEN      string                `json:"description_en"`
}

type PeriodicalIssueDTO struct {
	PeriodicalID int64  `json:"periodical_id"`
	IssueNumber  string `json:"issue_number"`
	IssueDate    string `json:"issue_date"` // YYYY-MM-DD
	AnnotationTM string `json:"annotation_tm"`
	AnnotationRU string `json:"annotation_ru"`
	AnnotationEN string `json:"annotation_en"`
}

// ---- Public ----

func (h *PeriodicalsHandler) List(c *fiber.Ctx) error {
	q := h.db.Model(&models.Periodical{}).Order("name_tm ASC")
	if t := c.Query("type"); t != "" {
		q = q.Where("type = ?", t)
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
	var items []models.Periodical
	q.Offset((page - 1) * limit).Limit(limit).Find(&items)
	return c.JSON(fiber.Map{"data": items, "total": total, "page": page, "limit": limit})
}

func (h *PeriodicalsHandler) Get(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid id"})
	}
	var p models.Periodical
	if err := h.db.First(&p, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "not found"})
	}
	return c.JSON(p)
}

func (h *PeriodicalsHandler) GetBySlug(c *fiber.Ctx) error {
	var p models.Periodical
	if err := h.db.Where("slug = ?", c.Params("slug")).First(&p).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "not found"})
	}
	return c.JSON(p)
}

// ListAllIssues — лента выпусков по всем изданиям с фильтрами.
// GET /api/periodicals/issues?type=&periodical_id=&before=YYYY-MM-DD&page=&limit=
// По умолчанию (без before) — самые свежие выпуски (issue_date DESC).
// before — показать выпуски с датой <= before (архив «на дату»).
func (h *PeriodicalsHandler) ListAllIssues(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "24"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 24
	}

	q := h.db.Model(&models.PeriodicalIssue{}).
		Joins("JOIN periodicals p ON p.id = periodical_issues.periodical_id")

	if t := c.Query("type"); t == string(models.PeriodicalNewspaper) || t == string(models.PeriodicalMagazine) {
		q = q.Where("p.type = ?", t)
	}
	if pid := c.Query("periodical_id"); pid != "" {
		if id, err := strconv.ParseInt(pid, 10, 64); err == nil {
			q = q.Where("periodical_issues.periodical_id = ?", id)
		}
	}
	if before := c.Query("before"); before != "" {
		if d, err := time.Parse("2006-01-02", before); err == nil {
			q = q.Where("periodical_issues.issue_date <= ?", d)
		}
	}

	var total int64
	q.Count(&total)

	var issues []models.PeriodicalIssue
	q.Preload("Periodical").
		Order("periodical_issues.issue_date DESC, periodical_issues.id DESC").
		Offset((page-1)*limit).Limit(limit).
		Find(&issues)

	return c.JSON(fiber.Map{"data": issues, "total": total, "page": page, "limit": limit})
}

func (h *PeriodicalsHandler) ListIssues(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid id"})
	}
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	q := h.db.Model(&models.PeriodicalIssue{}).Where("periodical_id = ?", id)
	if y := c.Query("year"); y != "" {
		q = q.Where("EXTRACT(YEAR FROM issue_date) = ?", y)
	}
	var total int64
	q.Count(&total)
	var issues []models.PeriodicalIssue
	q.Order("issue_date DESC").
		Offset((page-1)*limit).Limit(limit).
		Find(&issues)
	return c.JSON(fiber.Map{"data": issues, "total": total, "page": page, "limit": limit})
}

// IssueYears — список годов выпусков со счётчиками для архива (GET /periodicals/:id/issue-years).
func (h *PeriodicalsHandler) IssueYears(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid id"})
	}
	type yearCount struct {
		Year  int   `json:"year"`
		Count int64 `json:"count"`
	}
	var rows []yearCount
	h.db.Model(&models.PeriodicalIssue{}).
		Select("EXTRACT(YEAR FROM issue_date)::int AS year, COUNT(*) AS count").
		Where("periodical_id = ?", id).
		Group("EXTRACT(YEAR FROM issue_date)").
		Order("year DESC").
		Scan(&rows)
	return c.JSON(fiber.Map{"data": rows})
}

func (h *PeriodicalsHandler) GetIssue(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid id"})
	}
	var issue models.PeriodicalIssue
	if err := h.db.First(&issue, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "not found"})
	}
	return c.JSON(issue)
}

func (h *PeriodicalsHandler) ServeIssueFile(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid id"})
	}
	var issue models.PeriodicalIssue
	if err := h.db.First(&issue, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "not found"})
	}
	if issue.FilePath == "" {
		return c.Status(404).JSON(fiber.Map{"error": "no file"})
	}
	h.db.Model(&issue).UpdateColumn("views_count", gorm.Expr("views_count + 1"))
	return c.SendFile(h.files.AbsPath(issue.FilePath), false)
}

func (h *PeriodicalsHandler) DownloadIssue(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid id"})
	}
	var issue models.PeriodicalIssue
	if err := h.db.First(&issue, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "not found"})
	}
	if issue.FilePath == "" {
		return c.Status(404).JSON(fiber.Map{"error": "no file"})
	}
	h.db.Model(&issue).UpdateColumn("downloads_count", gorm.Expr("downloads_count + 1"))
	c.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="issue-%d.pdf"`, issue.ID))
	return c.Download(h.files.AbsPath(issue.FilePath))
}

// ---- Admin ----

func (h *PeriodicalsHandler) Create(c *fiber.Ctx) error {
	var dto PeriodicalDTO
	if err := c.BodyParser(&dto); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	if dto.NameTM == "" {
		return c.Status(400).JSON(fiber.Map{"error": "name_tm required"})
	}
	if dto.Type != models.PeriodicalNewspaper && dto.Type != models.PeriodicalMagazine {
		return c.Status(400).JSON(fiber.Map{"error": "type must be newspaper or magazine"})
	}
	if dto.Slug == "" {
		dto.Slug = slugify(dto.NameTM)
	}
	p := models.Periodical{
		Slug:        dto.Slug,
		Type:        dto.Type,
		Periodicity: dto.Periodicity,
		NameTM:      dto.NameTM,
		NameRU:      dto.NameRU,
		NameEN:      dto.NameEN,
		DescTM:      dto.DescTM,
		DescRU:      dto.DescRU,
		DescEN:      dto.DescEN,
	}
	if err := h.db.Create(&p).Error; err != nil {
		return c.Status(409).JSON(fiber.Map{"error": "slug already exists"})
	}
	return c.Status(201).JSON(p)
}

func (h *PeriodicalsHandler) Update(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid id"})
	}
	var p models.Periodical
	if err := h.db.First(&p, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "not found"})
	}
	var dto PeriodicalDTO
	if err := c.BodyParser(&dto); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	if dto.Slug != "" {
		p.Slug = dto.Slug
	}
	if dto.Type == models.PeriodicalNewspaper || dto.Type == models.PeriodicalMagazine {
		p.Type = dto.Type
	}
	p.Periodicity = dto.Periodicity
	p.NameTM = dto.NameTM
	p.NameRU = dto.NameRU
	p.NameEN = dto.NameEN
	p.DescTM = dto.DescTM
	p.DescRU = dto.DescRU
	p.DescEN = dto.DescEN

	if err := h.db.Save(&p).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(p)
}

func (h *PeriodicalsHandler) Delete(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid id"})
	}
	var count int64
	h.db.Model(&models.PeriodicalIssue{}).Where("periodical_id = ?", id).Count(&count)
	if count > 0 {
		return c.Status(409).JSON(fiber.Map{"error": "has issues, delete them first"})
	}
	h.db.Delete(&models.Periodical{}, id)
	return c.SendStatus(204)
}

func (h *PeriodicalsHandler) UploadCover(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid id"})
	}
	var p models.Periodical
	if err := h.db.First(&p, id).Error; err != nil {
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
	rel := fmt.Sprintf("periodicals/%d/cover%s", id, ext)
	if err := h.files.SaveUpload(fh, rel); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	h.db.Model(&p).Update("cover_path", rel)
	return c.JSON(fiber.Map{"cover_path": rel})
}

// ---- Issue admin ----

func (h *PeriodicalsHandler) CreateIssue(c *fiber.Ctx) error {
	pid, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid periodical id"})
	}
	var p models.Periodical
	if err := h.db.First(&p, pid).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "periodical not found"})
	}
	var dto PeriodicalIssueDTO
	if err := c.BodyParser(&dto); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	issueDate, err := time.Parse("2006-01-02", dto.IssueDate)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid issue_date, use YYYY-MM-DD"})
	}
	issue := models.PeriodicalIssue{
		PeriodicalID: pid,
		IssueNumber:  dto.IssueNumber,
		IssueDate:    issueDate,
		AnnotationTM: dto.AnnotationTM,
		AnnotationRU: dto.AnnotationRU,
		AnnotationEN: dto.AnnotationEN,
	}
	if err := h.db.Create(&issue).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(201).JSON(issue)
}

func (h *PeriodicalsHandler) UpdateIssue(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid id"})
	}
	var issue models.PeriodicalIssue
	if err := h.db.First(&issue, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "not found"})
	}
	var dto PeriodicalIssueDTO
	if err := c.BodyParser(&dto); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	if dto.IssueDate != "" {
		issueDate, err := time.Parse("2006-01-02", dto.IssueDate)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "invalid issue_date"})
		}
		issue.IssueDate = issueDate
	}
	issue.IssueNumber = dto.IssueNumber
	issue.AnnotationTM = dto.AnnotationTM
	issue.AnnotationRU = dto.AnnotationRU
	issue.AnnotationEN = dto.AnnotationEN
	h.db.Save(&issue)
	return c.JSON(issue)
}

func (h *PeriodicalsHandler) DeleteIssue(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid id"})
	}
	h.db.Delete(&models.PeriodicalIssue{}, id)
	return c.SendStatus(204)
}

func (h *PeriodicalsHandler) UploadIssueFile(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid id"})
	}
	var issue models.PeriodicalIssue
	if err := h.db.First(&issue, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "not found"})
	}
	fh, err := c.FormFile("file")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "file required"})
	}
	if !isPDFFile(fh) {
		return c.Status(400).JSON(fiber.Map{"error": "only PDF allowed"})
	}
	rel := fmt.Sprintf("periodicals/%d/issues/%d/file.pdf", issue.PeriodicalID, id)
	if err := h.files.SaveUpload(fh, rel); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	h.db.Model(&issue).Update("file_path", rel)

	// Если обложка ещё не загружена вручную — генерируем из первой страницы PDF.
	if issue.CoverPath == "" {
		outDir := fmt.Sprintf("periodicals/%d/issues/%d", issue.PeriodicalID, id)
		if coverPath, err := h.files.GenerateCoverFromPDF(rel, outDir); err == nil {
			h.db.Model(&issue).Update("cover_path", coverPath)
		}
	}

	return c.JSON(fiber.Map{"file_path": rel})
}

func (h *PeriodicalsHandler) UploadIssueCover(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid id"})
	}
	var issue models.PeriodicalIssue
	if err := h.db.First(&issue, id).Error; err != nil {
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
	rel := fmt.Sprintf("periodicals/%d/issues/%d/cover%s", issue.PeriodicalID, id, ext)
	if err := h.files.SaveUpload(fh, rel); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	h.db.Model(&issue).Update("cover_path", rel)
	return c.JSON(fiber.Map{"cover_path": rel})
}

// helpers shared with books.go via same package
func isImageFile(fh *multipart.FileHeader) bool {
	f, err := fh.Open()
	if err != nil {
		return false
	}
	defer f.Close()
	buf := make([]byte, 8)
	n, _ := f.Read(buf)
	if n < 4 {
		return false
	}
	// JPEG
	if buf[0] == 0xFF && buf[1] == 0xD8 {
		return true
	}
	// PNG
	if buf[0] == 0x89 && buf[1] == 0x50 && buf[2] == 0x4E && buf[3] == 0x47 {
		return true
	}
	// GIF
	if buf[0] == 0x47 && buf[1] == 0x49 && buf[2] == 0x46 {
		return true
	}
	// WebP: RIFF????WEBP
	if n >= 8 && buf[0] == 0x52 && buf[1] == 0x49 && buf[4] == 0x57 && buf[5] == 0x45 {
		return true
	}
	return false
}

func isPDFFile(fh *multipart.FileHeader) bool {
	f, err := fh.Open()
	if err != nil {
		return false
	}
	defer f.Close()
	buf := make([]byte, 4)
	n, _ := f.Read(buf)
	return n == 4 && buf[0] == '%' && buf[1] == 'P' && buf[2] == 'D' && buf[3] == 'F'
}

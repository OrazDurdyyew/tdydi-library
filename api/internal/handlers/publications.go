package handlers

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"library/internal/models"
	"library/internal/services"
)

type PublicationsHandler struct {
	db    *gorm.DB
	files *services.FileStorage
}

func NewPublicationsHandler(db *gorm.DB, files *services.FileStorage) *PublicationsHandler {
	return &PublicationsHandler{db: db, files: files}
}

type PublicationDTO struct {
	Slug              string `json:"slug"`
	TeacherID         int64  `json:"teacher_id"`
	CategoryID        int64  `json:"category_id"`
	PeriodicalIssueID *int64 `json:"periodical_issue_id"`
	NameTM            string `json:"name_tm"`
	NameRU            string `json:"name_ru"`
	NameEN            string `json:"name_en"`
	ArticleDate       string `json:"article_date"` // YYYY-MM-DD
	PageCount         int    `json:"page_count"`
	IsPublished       bool   `json:"is_published"`
}

func (h *PublicationsHandler) List(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	q := h.db.Model(&models.Publication{})

	if c.Locals("admin_id") == nil {
		q = q.Where("publications.is_published = ?", true)
	}

	if cat := c.Query("category_id"); cat != "" {
		if id, err := strconv.ParseInt(cat, 10, 64); err == nil {
			q = q.Where("publications.category_id = ?", id)
		}
	}
	if teacher := c.Query("teacher_id"); teacher != "" {
		if id, err := strconv.ParseInt(teacher, 10, 64); err == nil {
			q = q.Where("publications.teacher_id = ?", id)
		}
	}
	if dep := c.Query("department_id"); dep != "" {
		if id, err := strconv.ParseInt(dep, 10, 64); err == nil {
			// У публикации нет department_id — выводим через преподавателя.
			q = q.Joins("JOIN teachers t ON t.id = publications.teacher_id").
				Where("t.department_id = ?", id)
		}
	}
	if per := c.Query("periodical_id"); per != "" {
		if id, err := strconv.ParseInt(per, 10, 64); err == nil {
			q = q.Joins("JOIN periodical_issues pi ON pi.id = publications.periodical_issue_id").
				Where("pi.periodical_id = ?", id)
		}
	}
	if y := c.Query("year"); y != "" {
		if year, err := strconv.Atoi(y); err == nil {
			q = q.Where("EXTRACT(YEAR FROM publications.article_date) = ?", year)
		}
	}
	if search := c.Query("q"); search != "" {
		like := "%" + search + "%"
		q = q.Where("publications.name_tm ILIKE ? OR publications.name_ru ILIKE ? OR publications.name_en ILIKE ?", like, like, like)
	}

	switch c.Query("sort", "-created_at") {
	case "name":
		q = q.Order("publications.name_tm ASC")
	case "-name":
		q = q.Order("publications.name_tm DESC")
	case "date":
		q = q.Order("publications.article_date ASC")
	case "-date":
		q = q.Order("publications.article_date DESC")
	case "views":
		q = q.Order("publications.views_count DESC")
	case "downloads":
		q = q.Order("publications.downloads_count DESC")
	case "created_at":
		q = q.Order("publications.created_at ASC")
	default:
		q = q.Order("publications.created_at DESC")
	}

	var total int64
	q.Count(&total)

	var pubs []models.Publication
	if err := q.Preload("Teacher").Preload("Teacher.Department").Preload("Category").Preload("PeriodicalIssue").
		Limit(limit).Offset((page - 1) * limit).Find(&pubs).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"data": pubs,
		"meta": fiber.Map{"total": total, "page": page, "limit": limit},
	})
}

func (h *PublicationsHandler) Get(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid id"})
	}
	var p models.Publication
	if err := h.db.Preload("Teacher").Preload("Teacher.Department").Preload("Category").Preload("PeriodicalIssue").First(&p, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "not found"})
	}
	if !p.IsPublished && c.Locals("admin_id") == nil {
		return c.Status(404).JSON(fiber.Map{"error": "not found"})
	}
	return c.JSON(p)
}

func (h *PublicationsHandler) GetBySlug(c *fiber.Ctx) error {
	var p models.Publication
	if err := h.db.Preload("Teacher").Preload("Teacher.Department").Preload("Category").Preload("PeriodicalIssue").
		Where("slug = ?", c.Params("slug")).First(&p).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "not found"})
	}
	if !p.IsPublished && c.Locals("admin_id") == nil {
		return c.Status(404).JSON(fiber.Map{"error": "not found"})
	}
	return c.JSON(p)
}

// GetByTeacher — GET /api/publications/teacher/:slug. Преподаватель + его публикации.
func (h *PublicationsHandler) GetByTeacher(c *fiber.Ctx) error {
	var t models.Teacher
	if err := h.db.Preload("Department").Where("slug = ?", c.Params("slug")).First(&t).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "not found"})
	}
	q := h.db.Model(&models.Publication{}).Where("teacher_id = ?", t.ID)
	if c.Locals("admin_id") == nil {
		q = q.Where("is_published = ?", true)
	}
	var pubs []models.Publication
	q.Preload("Category").Preload("PeriodicalIssue").
		Order("article_date DESC, created_at DESC").Find(&pubs)
	return c.JSON(fiber.Map{"teacher": t, "publications": pubs})
}

func (h *PublicationsHandler) Create(c *fiber.Ctx) error {
	var dto PublicationDTO
	if err := c.BodyParser(&dto); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	if dto.NameTM == "" {
		return c.Status(400).JSON(fiber.Map{"error": "name_tm required"})
	}
	if dto.TeacherID == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "teacher_id required"})
	}
	if dto.CategoryID == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "category_id required"})
	}
	slug := dto.Slug
	if slug == "" {
		slug = slugify(dto.NameTM)
	}

	p := models.Publication{
		Slug:              slug,
		TeacherID:         dto.TeacherID,
		CategoryID:        dto.CategoryID,
		PeriodicalIssueID: dto.PeriodicalIssueID,
		NameTM:            dto.NameTM,
		NameRU:            dto.NameRU,
		NameEN:            dto.NameEN,
		PageCount:         dto.PageCount,
		IsPublished:       dto.IsPublished,
	}
	if dto.ArticleDate != "" {
		d, err := time.Parse("2006-01-02", dto.ArticleDate)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "invalid article_date, use YYYY-MM-DD"})
		}
		p.ArticleDate = &d
	}

	if err := h.db.Create(&p).Error; err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	h.db.Preload("Teacher").Preload("Category").Preload("PeriodicalIssue").First(&p, p.ID)
	return c.Status(201).JSON(p)
}

func (h *PublicationsHandler) Update(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid id"})
	}
	var p models.Publication
	if err := h.db.First(&p, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "not found"})
	}
	var dto PublicationDTO
	if err := c.BodyParser(&dto); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	if dto.Slug != "" {
		p.Slug = dto.Slug
	}
	if dto.TeacherID != 0 {
		p.TeacherID = dto.TeacherID
	}
	if dto.CategoryID != 0 {
		p.CategoryID = dto.CategoryID
	}
	p.PeriodicalIssueID = dto.PeriodicalIssueID
	if dto.NameTM != "" {
		p.NameTM = dto.NameTM
	}
	p.NameRU = dto.NameRU
	p.NameEN = dto.NameEN
	p.PageCount = dto.PageCount
	p.IsPublished = dto.IsPublished
	if dto.ArticleDate != "" {
		d, err := time.Parse("2006-01-02", dto.ArticleDate)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "invalid article_date"})
		}
		p.ArticleDate = &d
	} else {
		p.ArticleDate = nil
	}

	if err := h.db.Save(&p).Error; err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	h.db.Preload("Teacher").Preload("Category").Preload("PeriodicalIssue").First(&p, p.ID)
	return c.JSON(p)
}

func (h *PublicationsHandler) Delete(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid id"})
	}
	var p models.Publication
	if err := h.db.First(&p, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "not found"})
	}
	if p.FilePath != "" {
		_ = h.files.Remove(p.FilePath)
	}
	if p.CoverPath != "" {
		_ = h.files.Remove(p.CoverPath)
	}
	h.db.Delete(&models.Publication{}, id)
	return c.SendStatus(204)
}

// UploadFile — POST /api/admin/publications/:id/file (multipart, field "file").
// Сохраняет PDF в publications/{id}/file.pdf и при отсутствии обложки генерирует её.
func (h *PublicationsHandler) UploadFile(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid id"})
	}
	var p models.Publication
	if err := h.db.First(&p, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "not found"})
	}
	fh, err := c.FormFile("file")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "file field is required"})
	}
	if !isPDFFile(fh) {
		return c.Status(400).JSON(fiber.Map{"error": "only PDF allowed"})
	}
	rel := fmt.Sprintf("publications/%d/file.pdf", id)
	if err := h.files.SaveUpload(fh, rel); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	h.db.Model(&p).Update("file_path", rel)

	if p.CoverPath == "" {
		outDir := fmt.Sprintf("publications/%d", id)
		if coverPath, err := h.files.GenerateCoverFromPDF(rel, outDir); err == nil {
			h.db.Model(&p).Update("cover_path", coverPath)
		}
	}

	h.db.First(&p, id)
	return c.JSON(p)
}

func (h *PublicationsHandler) ServeFile(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid id"})
	}
	var p models.Publication
	if err := h.db.First(&p, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "not found"})
	}
	if !p.IsPublished {
		return c.Status(404).JSON(fiber.Map{"error": "not found"})
	}
	if p.FilePath == "" {
		return c.Status(404).JSON(fiber.Map{"error": "no file"})
	}
	h.db.Model(&p).UpdateColumn("views_count", gorm.Expr("views_count + 1"))
	return c.SendFile(h.files.AbsPath(p.FilePath), false)
}

func (h *PublicationsHandler) Download(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid id"})
	}
	var p models.Publication
	if err := h.db.First(&p, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "not found"})
	}
	if !p.IsPublished {
		return c.Status(404).JSON(fiber.Map{"error": "not found"})
	}
	if p.FilePath == "" {
		return c.Status(404).JSON(fiber.Map{"error": "no file"})
	}
	h.db.Model(&p).UpdateColumn("downloads_count", gorm.Expr("downloads_count + 1"))
	c.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="publication-%d.pdf"`, p.ID))
	return c.Download(h.files.AbsPath(p.FilePath))
}

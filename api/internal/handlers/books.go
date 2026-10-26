package handlers

import (
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"library/internal/models"
	"library/internal/services"
)

type BooksHandler struct {
	db    *gorm.DB
	perms *services.PermissionsService
	files *services.FileStorage
}

func NewBooksHandler(db *gorm.DB, perms *services.PermissionsService, files *services.FileStorage) *BooksHandler {
	return &BooksHandler{db: db, perms: perms, files: files}
}

type BookDTO struct {
	Slug            string  `json:"slug"`
	ISBN            string  `json:"isbn"`
	Year            int16   `json:"year"`
	Pages           int     `json:"pages"`
	Language        string  `json:"language"`
	Title           string  `json:"title"`
	Annotation      string  `json:"annotation"`
	Authors         string  `json:"authors"`
	Publisher       string  `json:"publisher"`
	AllowDownload   bool    `json:"allow_download"`
	IsPublished     bool    `json:"is_published"`
	IsInstituteBook bool    `json:"is_institute_book"`
	SectionIDs      []int64 `json:"section_ids"`
}

func (h *BooksHandler) List(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	q := h.db.Model(&models.Book{})

	if c.Locals("admin_id") == nil {
		q = q.Where("is_published = ?", true)
	} else {
		// Контент-менеджер в админке видит только книги из назначенных ему разделов.
		// Для суперадмина/all_sections AllowedSectionIDs возвращает nil — ограничений нет.
		adminID, _ := c.Locals("admin_id").(int64)
		allowed, err := h.perms.AllowedSectionIDs(adminID)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		if allowed != nil {
			if len(allowed) == 0 {
				q = q.Where("1 = 0")
			} else {
				q = q.Where("books.id IN (?)",
					h.db.Model(&models.BookSection{}).Select("book_id").Where("section_id IN ?", allowed))
			}
		}
	}

	if sectionIDStr := c.Query("section_id"); sectionIDStr != "" {
		if sid, err := strconv.ParseInt(sectionIDStr, 10, 64); err == nil {
			q = q.Joins("JOIN book_sections bs ON bs.book_id = books.id").
				Where("bs.section_id = ?", sid)
		}
	}
	if y := c.Query("year"); y != "" {
		if year, err := strconv.Atoi(y); err == nil {
			q = q.Where("year = ?", year)
		}
	}
	if lang := c.Query("language"); lang != "" {
		q = q.Where("language = ?", lang)
	}
	if author := c.Query("author"); author != "" {
		q = q.Where("authors ILIKE ?", "%"+author+"%")
	}
	if search := c.Query("q"); search != "" {
		q = q.Where("search_vector @@ websearch_to_tsquery('simple', unaccent(?))", search)
	}

	switch c.Query("sort", "-created_at") {
	case "title":
		q = q.Order("title ASC")
	case "-title":
		q = q.Order("title DESC")
	case "year":
		q = q.Order("year ASC")
	case "-year":
		q = q.Order("year DESC")
	case "views":
		q = q.Order("views_count DESC")
	case "downloads":
		q = q.Order("downloads_count DESC")
	case "created_at":
		q = q.Order("created_at ASC")
	default:
		q = q.Order("created_at DESC")
	}

	var total int64
	q.Count(&total)

	var books []models.Book
	if err := q.Preload("Sections").Limit(limit).Offset((page - 1) * limit).Find(&books).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"data": books,
		"meta": fiber.Map{"total": total, "page": page, "limit": limit},
	})
}

func (h *BooksHandler) Get(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid id"})
	}
	var book models.Book
	if err := h.db.Preload("Sections").Preload("Tags").First(&book, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "not found"})
	}
	if !book.IsPublished && c.Locals("admin_id") == nil {
		return c.Status(404).JSON(fiber.Map{"error": "not found"})
	}
	return c.JSON(book)
}

func (h *BooksHandler) GetBySlug(c *fiber.Ctx) error {
	slug := c.Params("slug")
	var book models.Book
	if err := h.db.Preload("Sections").Preload("Tags").Where("slug = ?", slug).First(&book).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "not found"})
	}
	if !book.IsPublished && c.Locals("admin_id") == nil {
		return c.Status(404).JSON(fiber.Map{"error": "not found"})
	}
	return c.JSON(book)
}

// Similar — GET /api/books/:id/similar. Книги, делящие хотя бы один раздел с указанной,
// отсортированные по популярности. Только опубликованные, без самой книги.
func (h *BooksHandler) Similar(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid id"})
	}
	var sectionIDs []int64
	h.db.Model(&models.BookSection{}).Where("book_id = ?", id).Pluck("section_id", &sectionIDs)

	books := []models.Book{}
	if len(sectionIDs) > 0 {
		h.db.Distinct("books.*").
			Joins("JOIN book_sections bs ON bs.book_id = books.id").
			Where("bs.section_id IN ? AND books.id <> ? AND books.is_published = ?", sectionIDs, id, true).
			Preload("Sections").
			Order("books.views_count DESC, books.created_at DESC").
			Limit(6).Find(&books)
	}
	return c.JSON(fiber.Map{"data": books})
}

func (h *BooksHandler) Create(c *fiber.Ctx) error {
	var dto BookDTO
	if err := c.BodyParser(&dto); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
	}
	if dto.Slug == "" || dto.Title == "" {
		return c.Status(400).JSON(fiber.Map{"error": "slug and title are required"})
	}
	if len(dto.SectionIDs) == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "at least one section is required"})
	}

	adminID, _ := c.Locals("admin_id").(int64)
	ok, err := h.perms.CanManageSections(adminID, dto.SectionIDs)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	if !ok {
		return c.Status(403).JSON(fiber.Map{"error": "no permission for specified sections"})
	}

	if dto.Language == "" {
		dto.Language = "tm"
	}

	book := models.Book{
		Slug:             dto.Slug,
		ISBN:             dto.ISBN,
		Year:             dto.Year,
		Pages:            dto.Pages,
		Language:         dto.Language,
		Title:            dto.Title,
		Annotation:       dto.Annotation,
		Authors:          dto.Authors,
		Publisher:        dto.Publisher,
		AllowDownload:    dto.AllowDownload,
		IsPublished:      dto.IsPublished,
		IsInstituteBook:  dto.IsInstituteBook,
		CreatedByAdminID: &adminID,
	}

	if err := h.db.Create(&book).Error; err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	for _, sid := range dto.SectionIDs {
		if err := h.db.Create(&models.BookSection{BookID: book.ID, SectionID: sid}).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
	}

	h.db.Preload("Sections").First(&book, book.ID)
	return c.Status(201).JSON(book)
}

func (h *BooksHandler) Update(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid id"})
	}

	var dto BookDTO
	if err := c.BodyParser(&dto); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
	}

	var existing models.Book
	if err := h.db.Preload("Sections").First(&existing, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "not found"})
	}

	adminID, _ := c.Locals("admin_id").(int64)

	currentSectionIDs := make([]int64, len(existing.Sections))
	for i, s := range existing.Sections {
		currentSectionIDs[i] = s.ID
	}
	if ok, err := h.perms.CanManageSections(adminID, currentSectionIDs); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	} else if !ok {
		return c.Status(403).JSON(fiber.Map{"error": "no permission for existing sections"})
	}

	if len(dto.SectionIDs) > 0 {
		if ok, err := h.perms.CanManageSections(adminID, dto.SectionIDs); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		} else if !ok {
			return c.Status(403).JSON(fiber.Map{"error": "no permission for new sections"})
		}
	}

	updates := map[string]interface{}{
		"slug":              dto.Slug,
		"isbn":              dto.ISBN,
		"year":              dto.Year,
		"pages":             dto.Pages,
		"language":          dto.Language,
		"title":             dto.Title,
		"annotation":        dto.Annotation,
		"authors":           dto.Authors,
		"publisher":         dto.Publisher,
		"allow_download":    dto.AllowDownload,
		"is_published":      dto.IsPublished,
		"is_institute_book": dto.IsInstituteBook,
	}
	if err := h.db.Model(&models.Book{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	if len(dto.SectionIDs) > 0 {
		h.db.Where("book_id = ?", id).Delete(&models.BookSection{})
		for _, sid := range dto.SectionIDs {
			h.db.Create(&models.BookSection{BookID: id, SectionID: sid})
		}
	}

	var book models.Book
	h.db.Preload("Sections").First(&book, id)
	return c.JSON(book)
}

func (h *BooksHandler) Delete(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid id"})
	}

	var existing models.Book
	if err := h.db.Preload("Sections").First(&existing, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "not found"})
	}

	adminID, _ := c.Locals("admin_id").(int64)
	sectionIDs := make([]int64, len(existing.Sections))
	for i, s := range existing.Sections {
		sectionIDs[i] = s.ID
	}
	if ok, err := h.perms.CanManageSections(adminID, sectionIDs); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	} else if !ok {
		return c.Status(403).JSON(fiber.Map{"error": "no permission to delete this book"})
	}

	h.db.Where("book_id = ?", id).Delete(&models.BookSection{})
	h.db.Where("book_id = ?", id).Delete(&models.BookTag{})

	// Удаляем файлы с диска (best-effort)
	if existing.FilePath != "" {
		_ = h.files.Remove(existing.FilePath)
	}
	if existing.CoverPath != "" {
		_ = h.files.Remove(existing.CoverPath)
	}

	if err := h.db.Delete(&models.Book{}, id).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"status": "deleted"})
}

// UploadFile — POST /api/admin/books/:id/file (multipart/form-data, field "file").
// Сохраняет файл в /var/library/books/{id}/file.{pdf|epub} и обновляет book.file_path.
// Для PDF автоматически генерирует обложку через pdftoppm (если её ещё нет).
func (h *BooksHandler) UploadFile(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid id"})
	}

	var book models.Book
	if err := h.db.Preload("Sections").First(&book, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "not found"})
	}

	adminID, _ := c.Locals("admin_id").(int64)
	sectionIDs := make([]int64, len(book.Sections))
	for i, s := range book.Sections {
		sectionIDs[i] = s.ID
	}
	if ok, err := h.perms.CanManageSections(adminID, sectionIDs); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	} else if !ok {
		return c.Status(403).JSON(fiber.Map{"error": "no permission"})
	}

	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "file field is required"})
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".pdf" && ext != ".epub" {
		return c.Status(400).JSON(fiber.Map{"error": "only pdf and epub are supported"})
	}

	bookDir := h.files.BookDir(id)
	if err := h.files.EnsureDir(bookDir); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	relPath := filepath.Join(bookDir, "file"+ext)
	if err := c.SaveFile(file, h.files.AbsPath(relPath)); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	fileFormat := strings.TrimPrefix(ext, ".")
	h.db.Model(&book).Updates(map[string]interface{}{
		"file_path":   relPath,
		"file_format": fileFormat,
	})

	if ext == ".pdf" && book.CoverPath == "" {
		if coverPath, err := h.files.GenerateCoverFromPDF(relPath, bookDir); err == nil {
			h.db.Model(&book).Update("cover_path", coverPath)
		}
	}

	h.db.First(&book, id)
	return c.JSON(book)
}

// UploadCover — POST /api/admin/books/:id/cover (multipart/form-data, field "cover").
// Сохраняет картинку в /var/library/books/{id}/cover.{ext} и пишет путь в book.cover_path.
// Если обложка задана вручную — автогенерация в UploadFile её не перезаписывает.
func (h *BooksHandler) UploadCover(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid id"})
	}

	var book models.Book
	if err := h.db.Preload("Sections").First(&book, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "not found"})
	}

	adminID, _ := c.Locals("admin_id").(int64)
	sectionIDs := make([]int64, len(book.Sections))
	for i, s := range book.Sections {
		sectionIDs[i] = s.ID
	}
	if ok, err := h.perms.CanManageSections(adminID, sectionIDs); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	} else if !ok {
		return c.Status(403).JSON(fiber.Map{"error": "no permission"})
	}

	fh, err := c.FormFile("cover")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "cover file required"})
	}
	if !isImageFile(fh) {
		return c.Status(400).JSON(fiber.Map{"error": "only images allowed"})
	}

	bookDir := h.files.BookDir(id)
	if err := h.files.EnsureDir(bookDir); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	ext := strings.ToLower(filepath.Ext(fh.Filename))
	rel := filepath.Join(bookDir, "cover"+ext)
	if err := h.files.SaveUpload(fh, rel); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	if book.CoverPath != "" && book.CoverPath != rel {
		_ = h.files.Remove(book.CoverPath)
	}

	h.db.Model(&book).Update("cover_path", rel)
	return c.JSON(fiber.Map{"cover_path": rel})
}

// ServeFile — GET /api/books/:id/file. Стримит файл книги для онлайн-чтения в браузере.
// Инкрементит views_count. Не проверяет allow_download (просмотр всегда разрешён, если книга опубликована).
func (h *BooksHandler) ServeFile(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid id"})
	}
	var book models.Book
	if err := h.db.First(&book, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "not found"})
	}
	if !book.IsPublished {
		return c.Status(404).JSON(fiber.Map{"error": "not found"})
	}
	if book.FilePath == "" {
		return c.Status(404).JSON(fiber.Map{"error": "no file"})
	}
	h.db.Model(&book).UpdateColumn("views_count", gorm.Expr("views_count + 1"))
	return c.SendFile(h.files.AbsPath(book.FilePath), false)
}

// Download — GET /api/books/:id/download. Скачивание файла.
// Требует allow_download=true. Инкрементит downloads_count.
func (h *BooksHandler) Download(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid id"})
	}
	var book models.Book
	if err := h.db.First(&book, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "not found"})
	}
	if !book.IsPublished {
		return c.Status(404).JSON(fiber.Map{"error": "not found"})
	}
	if !book.AllowDownload {
		return c.Status(403).JSON(fiber.Map{"error": "download not allowed"})
	}
	if book.FilePath == "" {
		return c.Status(404).JSON(fiber.Map{"error": "no file"})
	}
	h.db.Model(&book).Clauses(clause.Returning{}).UpdateColumn("downloads_count", gorm.Expr("downloads_count + 1"))
	return c.Download(h.files.AbsPath(book.FilePath))
}

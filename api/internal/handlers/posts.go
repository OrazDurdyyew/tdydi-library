package handlers

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"library/internal/models"
	"library/internal/services"
)

type PostsHandler struct {
	db    *gorm.DB
	files *services.FileStorage
}

func NewPostsHandler(db *gorm.DB, files *services.FileStorage) *PostsHandler {
	return &PostsHandler{db: db, files: files}
}

type PostDTO struct {
	Slug        string             `json:"slug"`
	Author      string             `json:"author"`
	Status      models.PostStatus  `json:"status"`
	PublishedAt string             `json:"published_at"` // RFC3339 or empty
	TitleTM     string             `json:"title_tm"`
	TitleRU     string             `json:"title_ru"`
	TitleEN     string             `json:"title_en"`
	ContentTM   string             `json:"content_tm"`
	ContentRU   string             `json:"content_ru"`
	ContentEN   string             `json:"content_en"`
}

// ---- Public ----

func (h *PostsHandler) List(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	q := h.db.Model(&models.Post{})

	isAdmin := c.Locals("admin_id") != nil
	if !isAdmin {
		q = q.Where("status = ?", models.PostStatusPublished)
	}

	var total int64
	q.Count(&total)
	var posts []models.Post
	q.Select("id, slug, cover_path, author, published_at, status, title_tm, title_ru, title_en, created_at, updated_at").
		Order("published_at DESC, created_at DESC").
		Offset((page - 1) * limit).Limit(limit).
		Find(&posts)
	return c.JSON(fiber.Map{"data": posts, "total": total, "page": page, "limit": limit})
}

func (h *PostsHandler) Get(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid id"})
	}
	var post models.Post
	if err := h.db.First(&post, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "not found"})
	}
	isAdmin := c.Locals("admin_id") != nil
	if !isAdmin && post.Status != models.PostStatusPublished {
		return c.Status(404).JSON(fiber.Map{"error": "not found"})
	}
	return c.JSON(post)
}

func (h *PostsHandler) GetBySlug(c *fiber.Ctx) error {
	var post models.Post
	if err := h.db.Where("slug = ?", c.Params("slug")).First(&post).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "not found"})
	}
	isAdmin := c.Locals("admin_id") != nil
	if !isAdmin && post.Status != models.PostStatusPublished {
		return c.Status(404).JSON(fiber.Map{"error": "not found"})
	}
	return c.JSON(post)
}

// ---- Admin ----

func (h *PostsHandler) Create(c *fiber.Ctx) error {
	var dto PostDTO
	if err := c.BodyParser(&dto); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	if dto.TitleTM == "" {
		return c.Status(400).JSON(fiber.Map{"error": "title_tm required"})
	}
	if dto.Slug == "" {
		dto.Slug = slugify(dto.TitleTM)
	}
	if dto.Status != models.PostStatusDraft && dto.Status != models.PostStatusPublished {
		dto.Status = models.PostStatusDraft
	}

	post := models.Post{
		Slug:      dto.Slug,
		Author:    dto.Author,
		Status:    dto.Status,
		TitleTM:   dto.TitleTM,
		TitleRU:   dto.TitleRU,
		TitleEN:   dto.TitleEN,
		ContentTM: dto.ContentTM,
		ContentRU: dto.ContentRU,
		ContentEN: dto.ContentEN,
	}
	if dto.PublishedAt != "" {
		t, err := time.Parse(time.RFC3339, dto.PublishedAt)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "invalid published_at"})
		}
		post.PublishedAt = &t
	} else if dto.Status == models.PostStatusPublished {
		now := time.Now()
		post.PublishedAt = &now
	}

	if err := h.db.Create(&post).Error; err != nil {
		return c.Status(409).JSON(fiber.Map{"error": "slug already exists"})
	}
	return c.Status(201).JSON(post)
}

func (h *PostsHandler) Update(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid id"})
	}
	var post models.Post
	if err := h.db.First(&post, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "not found"})
	}
	var dto PostDTO
	if err := c.BodyParser(&dto); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	if dto.Slug != "" {
		post.Slug = dto.Slug
	}
	post.Author = dto.Author
	if dto.Status == models.PostStatusDraft || dto.Status == models.PostStatusPublished {
		if post.Status != models.PostStatusPublished && dto.Status == models.PostStatusPublished && post.PublishedAt == nil {
			now := time.Now()
			post.PublishedAt = &now
		}
		post.Status = dto.Status
	}
	if dto.PublishedAt != "" {
		t, err := time.Parse(time.RFC3339, dto.PublishedAt)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "invalid published_at"})
		}
		post.PublishedAt = &t
	}
	post.TitleTM = dto.TitleTM
	post.TitleRU = dto.TitleRU
	post.TitleEN = dto.TitleEN
	post.ContentTM = dto.ContentTM
	post.ContentRU = dto.ContentRU
	post.ContentEN = dto.ContentEN

	h.db.Save(&post)
	return c.JSON(post)
}

func (h *PostsHandler) Delete(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid id"})
	}
	h.db.Delete(&models.Post{}, id)
	return c.SendStatus(204)
}

func (h *PostsHandler) UploadCover(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid id"})
	}
	var post models.Post
	if err := h.db.First(&post, id).Error; err != nil {
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
	rel := fmt.Sprintf("posts/%d/cover%s", id, ext)
	if err := h.files.SaveUpload(fh, rel); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	h.db.Model(&post).Update("cover_path", rel)
	return c.JSON(fiber.Map{"cover_path": rel})
}

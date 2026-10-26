package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2"

	"library/internal/services"
)

// ContentUploadsHandler handles uploads of inline content images
// (used by the WYSIWYG editor when authoring news posts / static pages).
type ContentUploadsHandler struct {
	files *services.FileStorage
}

func NewContentUploadsHandler(files *services.FileStorage) *ContentUploadsHandler {
	return &ContentUploadsHandler{files: files}
}

// UploadImage — POST /api/admin/content/image
// multipart/form-data with field "file". Returns {"path": "...", "url": "/api/files/..."}.
func (h *ContentUploadsHandler) UploadImage(c *fiber.Ctx) error {
	fh, err := c.FormFile("file")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "file field is required"})
	}
	if !isImageFile(fh) {
		return c.Status(400).JSON(fiber.Map{"error": "only image files are allowed (png, jpg, gif, webp)"})
	}

	ext := strings.ToLower(filepath.Ext(fh.Filename))
	if ext == "" {
		ext = ".png"
	}

	id, err := randomID(16)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "rand failed"})
	}
	rel := fmt.Sprintf("content/images/%s%s", id, ext)
	if err := h.files.SaveUpload(fh, rel); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{
		"path": rel,
		"url":  "/api/files/" + rel,
	})
}

func randomID(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

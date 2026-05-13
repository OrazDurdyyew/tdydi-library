package handlers

import (
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type FilesHandler struct {
	storagePath string
}

func NewFilesHandler(storagePath string) *FilesHandler {
	return &FilesHandler{storagePath: storagePath}
}

// ServePublic стримит файл из STORAGE_PATH. Используется для публичных файлов:
// обложек, логотипов, фавиконов. Файлы самих книг (file.pdf / file.epub) сюда
// не отдаются — для них есть отдельный download/serve endpoint с проверкой прав.
//
// Защита от path traversal: filepath.Clean + проверка отсутствия ".."
func (h *FilesHandler) ServePublic(c *fiber.Ctx) error {
	rel := c.Params("*")
	if rel == "" {
		return c.Status(404).JSON(fiber.Map{"error": "not found"})
	}

	rel = filepath.Clean(rel)
	if strings.Contains(rel, "..") || strings.HasPrefix(rel, "/") || strings.HasPrefix(rel, "\\") {
		return c.Status(400).JSON(fiber.Map{"error": "invalid path"})
	}

	// Файлы книг и выпусков не отдаём через публичный endpoint
	base := filepath.Base(rel)
	if strings.HasPrefix(base, "file.") {
		return c.Status(403).JSON(fiber.Map{"error": "use download endpoint"})
	}

	abs := filepath.Join(h.storagePath, rel)
	return c.SendFile(abs, false)
}

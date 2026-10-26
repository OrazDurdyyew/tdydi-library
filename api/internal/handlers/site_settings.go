package handlers

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/datatypes"
	"gorm.io/gorm"

	"library/internal/models"
	"library/internal/services"
)

type SiteSettingsHandler struct {
	db    *gorm.DB
	files *services.FileStorage
}

func NewSiteSettingsHandler(db *gorm.DB, files *services.FileStorage) *SiteSettingsHandler {
	return &SiteSettingsHandler{db: db, files: files}
}

func (h *SiteSettingsHandler) Get(c *fiber.Ctx) error {
	var s models.SiteSettings
	if err := h.db.First(&s, 1).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(s)
}

// SiteSettingsDTO — редактируемые из формы поля. Пути к логотипу/фавикону
// сюда НЕ входят: они задаются только через UploadAsset, иначе Update затирал бы их.
type SiteSettingsDTO struct {
	NameTM          string         `json:"name_tm"`
	NameRU          string         `json:"name_ru"`
	NameEN          string         `json:"name_en"`
	PrimaryColor    string         `json:"primary_color"`
	SecondaryColor  string         `json:"secondary_color"`
	AccentColor     string         `json:"accent_color"`
	BackgroundColor string         `json:"background_color"`
	ForegroundColor string         `json:"foreground_color"`
	MutedColor      string         `json:"muted_color"`
	BorderColor     string         `json:"border_color"`
	RingColor       string         `json:"ring_color"`
	DarkPrimaryColor    string     `json:"dark_primary_color"`
	DarkSecondaryColor  string     `json:"dark_secondary_color"`
	DarkAccentColor     string     `json:"dark_accent_color"`
	DarkBackgroundColor string     `json:"dark_background_color"`
	DarkForegroundColor string     `json:"dark_foreground_color"`
	DarkMutedColor      string     `json:"dark_muted_color"`
	DarkBorderColor     string     `json:"dark_border_color"`
	DarkRingColor       string     `json:"dark_ring_color"`
	FontFamily      string         `json:"font_family"`
	FontURL         string         `json:"font_url"`
	ContactAddress  string         `json:"contact_address"`
	ContactPhone    string         `json:"contact_phone"`
	ContactEmail    string         `json:"contact_email"`
	SocialLinks     datatypes.JSON `json:"social_links"`
	FooterTM        string         `json:"footer_tm"`
	FooterRU        string         `json:"footer_ru"`
	FooterEN        string         `json:"footer_en"`
}

func (h *SiteSettingsHandler) Update(c *fiber.Ctx) error {
	var dto SiteSettingsDTO
	if err := c.BodyParser(&dto); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
	}

	social := dto.SocialLinks
	if len(social) == 0 {
		social = datatypes.JSON([]byte("{}"))
	}

	updates := map[string]interface{}{
		"name_tm":          dto.NameTM,
		"name_ru":          dto.NameRU,
		"name_en":          dto.NameEN,
		"primary_color":    dto.PrimaryColor,
		"secondary_color":  dto.SecondaryColor,
		"accent_color":     dto.AccentColor,
		"background_color": dto.BackgroundColor,
		"foreground_color": dto.ForegroundColor,
		"muted_color":      dto.MutedColor,
		"border_color":     dto.BorderColor,
		"ring_color":       dto.RingColor,
		"dark_primary_color":    dto.DarkPrimaryColor,
		"dark_secondary_color":  dto.DarkSecondaryColor,
		"dark_accent_color":     dto.DarkAccentColor,
		"dark_background_color": dto.DarkBackgroundColor,
		"dark_foreground_color": dto.DarkForegroundColor,
		"dark_muted_color":      dto.DarkMutedColor,
		"dark_border_color":     dto.DarkBorderColor,
		"dark_ring_color":       dto.DarkRingColor,
		"font_family":      dto.FontFamily,
		"font_url":         dto.FontURL,
		"contact_address":  dto.ContactAddress,
		"contact_phone":    dto.ContactPhone,
		"contact_email":    dto.ContactEmail,
		"social_links":     social,
		"footer_tm":        dto.FooterTM,
		"footer_ru":        dto.FooterRU,
		"footer_en":        dto.FooterEN,
	}
	if err := h.db.Model(&models.SiteSettings{}).Where("id = ?", 1).Updates(updates).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	var s models.SiteSettings
	h.db.First(&s, 1)
	return c.JSON(s)
}

// UploadAsset — POST /api/admin/site-settings/asset?kind=logo_light|logo_dark|favicon
// Сохраняет файл в branding/ и обновляет соответствующий путь в site_settings.
func (h *SiteSettingsHandler) UploadAsset(c *fiber.Ctx) error {
	var col string
	switch c.Query("kind") {
	case "logo_light":
		col = "logo_light_path"
	case "logo_dark":
		col = "logo_dark_path"
	case "favicon":
		col = "favicon_path"
	default:
		return c.Status(400).JSON(fiber.Map{"error": "kind must be logo_light, logo_dark or favicon"})
	}

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
	rel := fmt.Sprintf("branding/%s%s", c.Query("kind"), ext)
	if err := h.files.SaveUpload(fh, rel); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	if err := h.db.Model(&models.SiteSettings{}).Where("id = ?", 1).Update(col, rel).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"path": rel})
}

// UploadFont — POST /api/admin/site-settings/font
// Сохраняет файл шрифта в fonts/ и обновляет соответствующий путь в site_settings.
func (h *SiteSettingsHandler) UploadFont(c *fiber.Ctx) error {
	fh, err := c.FormFile("file")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "file field is required"})
	}

	// Проверяем расширение файла
	ext := strings.ToLower(filepath.Ext(fh.Filename))
	allowed := []string{".ttf", ".otf", ".woff", ".woff2", ".eot"}
	valid := false
	for _, a := range allowed {
		if ext == a {
			valid = true
			break
		}
	}
	if !valid {
		return c.Status(400).JSON(fiber.Map{"error": "only font files are allowed (ttf, otf, woff, woff2, eot)"})
	}

	// Генерируем уникальное имя файла
	fontName := strings.TrimSuffix(filepath.Base(fh.Filename), ext)
	rel := fmt.Sprintf("fonts/%s-%d%s", fontName, time.Now().Unix(), ext)
	
	if err := h.files.SaveUpload(fh, rel); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	// Обновляем запись в базе
	updates := map[string]interface{}{
		"font_file_path":   rel,
		"font_uploaded_at": time.Now(),
	}
	if err := h.db.Model(&models.SiteSettings{}).Where("id = ?", 1).Updates(updates).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	// Возвращаем путь к файлу
	return c.JSON(fiber.Map{"path": rel})
}

// DeleteFont — DELETE /api/admin/site-settings/font
// Удаляет загруженный файл шрифта и очищает поля в site_settings.
func (h *SiteSettingsHandler) DeleteFont(c *fiber.Ctx) error {
	var s models.SiteSettings
	if err := h.db.First(&s, 1).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	// Удаляем файл, если он существует
	if s.FontFilePath != "" {
		if err := h.files.Remove(s.FontFilePath); err != nil {
			// Не прерываем выполнение, если файл не найден
			fmt.Printf("Warning: could not delete font file %s: %v\n", s.FontFilePath, err)
		}
	}

	// Очищаем поля в базе
	updates := map[string]interface{}{
		"font_file_path":   nil,
		"font_uploaded_at": nil,
	}
	if err := h.db.Model(&models.SiteSettings{}).Where("id = ?", 1).Updates(updates).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "font deleted successfully"})
}
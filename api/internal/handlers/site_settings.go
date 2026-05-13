package handlers

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"library/internal/models"
)

type SiteSettingsHandler struct {
	db *gorm.DB
}

func NewSiteSettingsHandler(db *gorm.DB) *SiteSettingsHandler {
	return &SiteSettingsHandler{db: db}
}

func (h *SiteSettingsHandler) Get(c *fiber.Ctx) error {
	var s models.SiteSettings
	if err := h.db.First(&s, 1).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(s)
}

func (h *SiteSettingsHandler) Update(c *fiber.Ctx) error {
	var s models.SiteSettings
	if err := c.BodyParser(&s); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
	}
	s.ID = 1
	if err := h.db.Save(&s).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(s)
}

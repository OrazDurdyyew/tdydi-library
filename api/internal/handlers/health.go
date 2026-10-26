package handlers

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"library/internal/version"
)

type HealthHandler struct {
	db *gorm.DB
}

func NewHealthHandler(db *gorm.DB) *HealthHandler {
	return &HealthHandler{db: db}
}

func (h *HealthHandler) Get(c *fiber.Ctx) error {
	sqlDB, err := h.db.DB()
	if err != nil {
		return c.Status(503).JSON(fiber.Map{"status": "unhealthy", "version": version.Version, "db": "no connection"})
	}
	if err := sqlDB.Ping(); err != nil {
		return c.Status(503).JSON(fiber.Map{"status": "unhealthy", "version": version.Version, "db": err.Error()})
	}
	return c.JSON(fiber.Map{"status": "ok", "version": version.Version, "db": "connected"})
}

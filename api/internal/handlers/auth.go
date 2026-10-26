package handlers

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"library/internal/auth"
	"library/internal/models"
)

type AuthHandler struct {
	db *gorm.DB
	tm *auth.TokenManager
}

func NewAuthHandler(db *gorm.DB, tm *auth.TokenManager) *AuthHandler {
	return &AuthHandler{db: db, tm: tm}
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AdminInfo struct {
	ID                int64   `json:"id"`
	Username          string  `json:"username"`
	Email             string  `json:"email,omitempty"`
	Role              string  `json:"role"`
	AllSections       bool    `json:"all_sections"`
	AllowedSectionIDs []int64 `json:"allowed_section_ids"`
}

type LoginResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	Admin        AdminInfo `json:"admin"`
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
	}
	if req.Username == "" || req.Password == "" {
		return c.Status(400).JSON(fiber.Map{"error": "username and password required"})
	}

	var admin models.Admin
	if err := h.db.Where("username = ?", req.Username).First(&admin).Error; err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "invalid credentials"})
	}
	if !auth.CheckPassword(req.Password, admin.PasswordHash) {
		return c.Status(401).JSON(fiber.Map{"error": "invalid credentials"})
	}

	access, err := h.tm.GenerateAccess(admin.ID, admin.Username, string(admin.Role))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "token generation failed"})
	}
	refresh, err := h.tm.GenerateRefresh(admin.ID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "token generation failed"})
	}

	allowed, err := h.allowedSectionIDs(admin.ID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "permissions fetch failed"})
	}

	return c.JSON(LoginResponse{
		AccessToken:  access,
		RefreshToken: refresh,
		Admin: AdminInfo{
			ID:                admin.ID,
			Username:          admin.Username,
			Email:             admin.Email,
			Role:              string(admin.Role),
			AllSections:       admin.AllSections,
			AllowedSectionIDs: allowed,
		},
	})
}

func (h *AuthHandler) Me(c *fiber.Ctx) error {
	adminID, _ := c.Locals("admin_id").(int64)
	var admin models.Admin
	if err := h.db.First(&admin, adminID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "admin not found"})
	}
	allowed, _ := h.allowedSectionIDs(admin.ID)
	return c.JSON(AdminInfo{
		ID:                admin.ID,
		Username:          admin.Username,
		Email:             admin.Email,
		Role:              string(admin.Role),
		AllSections:       admin.AllSections,
		AllowedSectionIDs: allowed,
	})
}

func (h *AuthHandler) allowedSectionIDs(adminID int64) ([]int64, error) {
	var perms []models.AdminSectionPermission
	if err := h.db.Where("admin_id = ?", adminID).Find(&perms).Error; err != nil {
		return nil, err
	}
	ids := make([]int64, 0, len(perms))
	for _, p := range perms {
		ids = append(ids, p.SectionID)
	}
	return ids, nil
}

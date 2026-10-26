package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"library/internal/models"
)

type AdminsHandler struct {
	db *gorm.DB
}

func NewAdminsHandler(db *gorm.DB) *AdminsHandler {
	return &AdminsHandler{db: db}
}

type AdminCreateDTO struct {
	Username    string            `json:"username"`
	Password    string            `json:"password"`
	Email       string            `json:"email"`
	Role        models.AdminRole  `json:"role"`
	AllSections bool              `json:"all_sections"`
	SectionIDs  []int64           `json:"section_ids"`
}

type AdminUpdateDTO struct {
	Email       string            `json:"email"`
	Role        models.AdminRole  `json:"role"`
	AllSections bool              `json:"all_sections"`
	Password    string            `json:"password"` // optional: set to reset
}

func (h *AdminsHandler) List(c *fiber.Ctx) error {
	var admins []models.Admin
	if err := h.db.Preload("SectionPermissions").Find(&admins).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(admins)
}

func (h *AdminsHandler) Get(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid id"})
	}
	var admin models.Admin
	if err := h.db.Preload("SectionPermissions").First(&admin, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "not found"})
	}
	return c.JSON(admin)
}

func (h *AdminsHandler) Create(c *fiber.Ctx) error {
	var dto AdminCreateDTO
	if err := c.BodyParser(&dto); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	if dto.Username == "" || dto.Password == "" {
		return c.Status(400).JSON(fiber.Map{"error": "username and password required"})
	}
	if dto.Role != models.RoleSuperadmin && dto.Role != models.RoleContentManager && dto.Role != models.RolePublicationManager {
		dto.Role = models.RoleContentManager
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(dto.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to hash password"})
	}

	admin := models.Admin{
		Username:     dto.Username,
		PasswordHash: string(hash),
		Email:        dto.Email,
		Role:         dto.Role,
		AllSections:  dto.AllSections,
	}

	if err := h.db.Create(&admin).Error; err != nil {
		return c.Status(409).JSON(fiber.Map{"error": "username already exists"})
	}

	if !dto.AllSections && len(dto.SectionIDs) > 0 {
		var perms []models.AdminSectionPermission
		for _, sid := range dto.SectionIDs {
			perms = append(perms, models.AdminSectionPermission{AdminID: admin.ID, SectionID: sid})
		}
		h.db.CreateInBatches(perms, 50)
	}

	h.db.Preload("SectionPermissions").First(&admin, admin.ID)
	return c.Status(201).JSON(admin)
}

func (h *AdminsHandler) Update(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid id"})
	}

	// Prevent self-demotion
	callerID, _ := c.Locals("admin_id").(int64)
	if callerID == id {
		return c.Status(400).JSON(fiber.Map{"error": "cannot modify your own account"})
	}

	var admin models.Admin
	if err := h.db.First(&admin, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "not found"})
	}

	var dto AdminUpdateDTO
	if err := c.BodyParser(&dto); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	if dto.Role == models.RoleSuperadmin || dto.Role == models.RoleContentManager || dto.Role == models.RolePublicationManager {
		admin.Role = dto.Role
	}
	admin.Email = dto.Email
	admin.AllSections = dto.AllSections

	if dto.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(dto.Password), bcrypt.DefaultCost)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "failed to hash password"})
		}
		admin.PasswordHash = string(hash)
	}

	if err := h.db.Save(&admin).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(admin)
}

func (h *AdminsHandler) Delete(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid id"})
	}
	callerID, _ := c.Locals("admin_id").(int64)
	if callerID == id {
		return c.Status(400).JSON(fiber.Map{"error": "cannot delete yourself"})
	}
	if err := h.db.Delete(&models.Admin{}, id).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	h.db.Where("admin_id = ?", id).Delete(&models.AdminSectionPermission{})
	return c.SendStatus(204)
}

type SetPermissionsDTO struct {
	AllSections bool    `json:"all_sections"`
	SectionIDs  []int64 `json:"section_ids"`
}

func (h *AdminsHandler) SetPermissions(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid id"})
	}
	var admin models.Admin
	if err := h.db.First(&admin, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "not found"})
	}
	var dto SetPermissionsDTO
	if err := c.BodyParser(&dto); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	admin.AllSections = dto.AllSections
	h.db.Save(&admin)

	h.db.Where("admin_id = ?", id).Delete(&models.AdminSectionPermission{})
	if !dto.AllSections && len(dto.SectionIDs) > 0 {
		var perms []models.AdminSectionPermission
		for _, sid := range dto.SectionIDs {
			perms = append(perms, models.AdminSectionPermission{AdminID: id, SectionID: sid})
		}
		h.db.CreateInBatches(perms, 50)
	}

	h.db.Preload("SectionPermissions").First(&admin, id)
	return c.JSON(admin)
}

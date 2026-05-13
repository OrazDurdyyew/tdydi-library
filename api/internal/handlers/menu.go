package handlers

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"library/internal/models"
)

type MenuHandler struct {
	db *gorm.DB
}

func NewMenuHandler(db *gorm.DB) *MenuHandler {
	return &MenuHandler{db: db}
}

type MenuSection struct {
	ID       int64         `json:"id"`
	Slug     string        `json:"slug"`
	NameTM   string        `json:"name_tm"`
	NameRU   string        `json:"name_ru"`
	NameEN   string        `json:"name_en"`
	Children []MenuSection `json:"children,omitempty"`
}

type MenuPage struct {
	Slug    string `json:"slug"`
	TitleTM string `json:"title_tm"`
	TitleRU string `json:"title_ru"`
	TitleEN string `json:"title_en"`
}

type MenuResponse struct {
	Sections []MenuSection `json:"sections"`
	Pages    []MenuPage    `json:"pages"`
}

// Get строит публичное меню сайта: дерево разделов (2 уровня) + статичные страницы с in_menu=true.
func (h *MenuHandler) Get(c *fiber.Ctx) error {
	var sections []models.Section
	if err := h.db.Order("sort_order ASC, id ASC").Find(&sections).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	childrenByParent := make(map[int64][]MenuSection)
	roots := make([]MenuSection, 0)
	rootIdx := make(map[int64]int) // section_id -> index in roots

	for _, s := range sections {
		m := MenuSection{
			ID:     s.ID,
			Slug:   s.Slug,
			NameTM: s.NameTM,
			NameRU: s.NameRU,
			NameEN: s.NameEN,
		}
		if s.ParentID == nil {
			rootIdx[s.ID] = len(roots)
			roots = append(roots, m)
		} else {
			childrenByParent[*s.ParentID] = append(childrenByParent[*s.ParentID], m)
		}
	}

	for parentID, kids := range childrenByParent {
		if idx, ok := rootIdx[parentID]; ok {
			roots[idx].Children = kids
		}
	}

	var pages []models.Page
	h.db.Where("in_menu = ?", true).Order("menu_order ASC, id ASC").Find(&pages)

	menuPages := make([]MenuPage, len(pages))
	for i, p := range pages {
		menuPages[i] = MenuPage{
			Slug:    p.Slug,
			TitleTM: p.TitleTM,
			TitleRU: p.TitleRU,
			TitleEN: p.TitleEN,
		}
	}

	return c.JSON(MenuResponse{Sections: roots, Pages: menuPages})
}

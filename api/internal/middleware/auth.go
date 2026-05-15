package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"

	"library/internal/auth"
)

// RequireAuth извлекает JWT из Authorization-заголовка, валидирует и кладёт claims в Locals.
// Используется для всех маршрутов /api/admin/*.
func RequireAuth(tm *auth.TokenManager) fiber.Handler {
	return func(c *fiber.Ctx) error {
		header := c.Get("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			return c.Status(401).JSON(fiber.Map{"error": "missing or invalid auth header"})
		}
		tokenStr := strings.TrimPrefix(header, "Bearer ")
		claims, err := tm.Verify(tokenStr)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{"error": "invalid token"})
		}
		c.Locals("admin_id", claims.AdminID)
		c.Locals("admin_role", claims.Role)
		c.Locals("admin_username", claims.Username)
		return c.Next()
	}
}

// OptionalAuth разбирает JWT, если он есть, и кладёт claims в Locals.
// В отличие от RequireAuth — не отклоняет запрос без токена.
// Нужен на публичных GET-маршрутах: позволяет залогиненному админу
// видеть черновики (постов, книг) через те же эндпоинты, что и гости.
func OptionalAuth(tm *auth.TokenManager) fiber.Handler {
	return func(c *fiber.Ctx) error {
		header := c.Get("Authorization")
		if strings.HasPrefix(header, "Bearer ") {
			tokenStr := strings.TrimPrefix(header, "Bearer ")
			if claims, err := tm.Verify(tokenStr); err == nil {
				c.Locals("admin_id", claims.AdminID)
				c.Locals("admin_role", claims.Role)
				c.Locals("admin_username", claims.Username)
			}
		}
		return c.Next()
	}
}

// RequireSuperadmin блокирует доступ всем, кроме роли superadmin.
// Применяется после RequireAuth.
func RequireSuperadmin() fiber.Handler {
	return func(c *fiber.Ctx) error {
		role, _ := c.Locals("admin_role").(string)
		if role != "superadmin" {
			return c.Status(403).JSON(fiber.Map{"error": "superadmin only"})
		}
		return c.Next()
	}
}

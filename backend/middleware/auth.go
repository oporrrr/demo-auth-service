package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

// ExtractBearerToken ดึง token จาก Authorization header แล้วใส่ใน Locals
func ExtractBearerToken(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"statusCode": 401,
			"code":       "UNAUTHORIZED",
			"message":    "Missing or invalid Authorization header",
		})
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	c.Locals("userToken", token)
	return c.Next()
}

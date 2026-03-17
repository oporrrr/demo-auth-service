package handler

import (
	"demo-auth-center/client"

	"github.com/gofiber/fiber/v2"
)

type MeHandler struct {
	roleClient *client.RoleClient
}

func NewMeHandler(roleClient *client.RoleClient) *MeHandler {
	return &MeHandler{roleClient: roleClient}
}

// GET /api/v1/me/permissions
func (h *MeHandler) GetPermissions(c *fiber.Ctx) error {
	accountID := c.Locals("accountId").(string)
	perms, err := h.roleClient.GetPermissions(accountID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"code": "ERROR", "message": err.Error()})
	}
	return c.JSON(fiber.Map{"code": "SUCCESS", "data": perms})
}

// GET /api/v1/me/menus
func (h *MeHandler) GetMenus(c *fiber.Ctx) error {
	accountID := c.Locals("accountId").(string)
	menus, err := h.roleClient.GetMenus(accountID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"code": "ERROR", "message": err.Error()})
	}
	return c.JSON(fiber.Map{"code": "SUCCESS", "data": menus})
}

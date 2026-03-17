package handler

import (
	"demo-auth-center/model"
	"demo-auth-center/service"

	"github.com/gofiber/fiber/v2"
)

type AccountHandler struct {
	svc *service.AuthCenterService
}

func NewAccountHandler() *AccountHandler {
	return &AccountHandler{svc: service.GetAuthCenterService()}
}

func (h *AccountHandler) CheckExistValue(c *fiber.Ctx) error {
	var req model.CheckExistValueRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "invalid request body"})
	}

	body, statusCode, err := h.svc.CheckExistValue(req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": err.Error()})
	}

	return sendRawJSON(c, statusCode, body)
}

func (h *AccountHandler) GetInformation(c *fiber.Ctx) error {
	userToken := c.Locals("userToken").(string)

	body, statusCode, err := h.svc.GetAccountInformation(userToken)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": err.Error()})
	}

	return sendRawJSON(c, statusCode, body)
}

func (h *AccountHandler) LinkCIS(c *fiber.Ctx) error {
	var req model.LinkCISRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "invalid request body"})
	}

	userToken := c.Locals("userToken").(string)
	body, statusCode, err := h.svc.LinkCIS(req, userToken)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": err.Error()})
	}

	return sendRawJSON(c, statusCode, body)
}

func (h *AccountHandler) UpdateUsername(c *fiber.Ctx) error {
	var req model.UpdateUsernameRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "invalid request body"})
	}

	userToken := c.Locals("userToken").(string)
	body, statusCode, err := h.svc.UpdateUsername(req, userToken)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": err.Error()})
	}

	return sendRawJSON(c, statusCode, body)
}

func (h *AccountHandler) UpdateProfile(c *fiber.Ctx) error {
	var req model.UpdateProfileRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "invalid request body"})
	}

	userToken := c.Locals("userToken").(string)
	body, statusCode, err := h.svc.UpdateProfile(req, userToken)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": err.Error()})
	}

	return sendRawJSON(c, statusCode, body)
}

func (h *AccountHandler) GetCISNumber(c *fiber.Ctx) error {
	var req model.GetCISNumberRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "invalid request body"})
	}

	userToken := c.Locals("userToken").(string)
	body, statusCode, err := h.svc.GetCISNumber(req, userToken)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": err.Error()})
	}

	return sendRawJSON(c, statusCode, body)
}

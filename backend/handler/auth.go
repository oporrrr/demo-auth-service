package handler

import (
	"encoding/json"

	"demo-auth-center/model"
	"demo-auth-center/service"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	svc *service.UserService
}

func NewAuthHandler(userSvc *service.UserService) *AuthHandler {
	return &AuthHandler{svc: userSvc}
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req model.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"code": "BAD_REQUEST", "message": "invalid request body"})
	}

	if req.Email == "" && req.PhoneNumber == "" {
		return c.Status(400).JSON(fiber.Map{"code": "BAD_REQUEST", "message": "email or phoneNumber is required"})
	}

	if req.PhoneNumber != "" && req.CountryCode == "" {
		return c.Status(400).JSON(fiber.Map{"code": "BAD_REQUEST", "message": "countryCode is required when using phoneNumber"})
	}

	body, statusCode, err := h.svc.Register(req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"code": "INTERNAL_ERROR", "message": err.Error()})
	}

	return sendRawJSON(c, statusCode, body)
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req model.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "invalid request body"})
	}

	body, statusCode, err := h.svc.Login(req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": err.Error()})
	}

	return sendRawJSON(c, statusCode, body)
}

func (h *AuthHandler) LoginWithOTP(c *fiber.Ctx) error {
	var req model.LoginWithOTPRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "invalid request body"})
	}

	body, statusCode, err := h.svc.LoginWithOTP(req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": err.Error()})
	}

	return sendRawJSON(c, statusCode, body)
}

func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	var req model.RefreshTokenRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "invalid request body"})
	}

	body, statusCode, err := h.svc.RefreshToken(req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": err.Error()})
	}

	return sendRawJSON(c, statusCode, body)
}

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	var req model.LogoutRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "invalid request body"})
	}

	userToken := c.Locals("userToken").(string)
	body, statusCode, err := h.svc.Logout(req, userToken)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": err.Error()})
	}

	return sendRawJSON(c, statusCode, body)
}

func (h *AuthHandler) UpdatePassword(c *fiber.Ctx) error {
	var req model.UpdatePasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "invalid request body"})
	}

	userToken := c.Locals("userToken").(string)
	body, statusCode, err := h.svc.UpdatePassword(req, userToken)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": err.Error()})
	}

	return sendRawJSON(c, statusCode, body)
}

// ── helper ────────────────────────────────────────────────

func sendRawJSON(c *fiber.Ctx, statusCode int, body []byte) error {
	var result interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return c.Status(statusCode).Send(body)
	}
	return c.Status(statusCode).JSON(result)
}

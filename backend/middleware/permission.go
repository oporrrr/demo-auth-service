package middleware

import (
	"encoding/json"
	"log"
	"time"

	"demo-auth-center/cache"
	"demo-auth-center/client"
	"demo-auth-center/service"

	"github.com/gofiber/fiber/v2"
)

type accountInfoResp struct {
	Code string `json:"code"`
	Data struct {
		ID string `json:"id"`
	} `json:"data"`
}

// ResolveAccountID resolves token → accountId via cache then auth center.
// Must run after ExtractBearerToken.
func ResolveAccountID(tokenCache *cache.TokenCache) fiber.Handler {
	authSvc := service.GetAuthCenterService()
	return func(c *fiber.Ctx) error {
		token, ok := c.Locals("userToken").(string)
		if !ok || token == "" {
			return c.Status(401).JSON(fiber.Map{"code": "UNAUTHORIZED"})
		}

		// cache hit
		if accountID, ok := tokenCache.Get(token); ok {
			c.Locals("accountId", accountID)
			return c.Next()
		}

		// cache miss — call auth center
		body, _, err := authSvc.GetAccountInformation(token)
		if err != nil {
			log.Printf("warn: ResolveAccountID failed: %v", err)
			return c.Status(401).JSON(fiber.Map{"code": "UNAUTHORIZED"})
		}
		var info accountInfoResp
		if err := json.Unmarshal(body, &info); err != nil || info.Code != "SUCCESS" || info.Data.ID == "" {
			return c.Status(401).JSON(fiber.Map{"code": "UNAUTHORIZED"})
		}

		tokenCache.SetWithTTL(token, info.Data.ID, 15*time.Minute)
		c.Locals("accountId", info.Data.ID)
		return c.Next()
	}
}

// RequirePermission checks if the resolved accountId has resource:action permission.
// Must run after ResolveAccountID.
func RequirePermission(roleClient *client.RoleClient, resource, action string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		accountID, ok := c.Locals("accountId").(string)
		if !ok || accountID == "" {
			return c.Status(401).JSON(fiber.Map{"code": "UNAUTHORIZED"})
		}
		if !roleClient.Check(accountID, resource, action) {
			return c.Status(403).JSON(fiber.Map{"code": "FORBIDDEN", "message": "permission denied: " + resource + ":" + action})
		}
		return c.Next()
	}
}

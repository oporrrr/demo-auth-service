package main

import (
	"log"
	"time"

	"demo-auth-center/cache"
	"demo-auth-center/client"
	"demo-auth-center/config"
	"demo-auth-center/database"
	"demo-auth-center/handler"
	"demo-auth-center/middleware"
	"demo-auth-center/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	// Load config
	config.Load()

	// Connect to database
	database.Connect()

	roleClient := client.NewRoleClient()

	// Shared token cache (token → accountId)
	tokenCache := cache.NewTokenCache(15 * time.Minute)

	// Init Fiber
	app := fiber.New(fiber.Config{
		AppName: "Demo Auth Center",
	})

	// Global middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:5173",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}))
	app.Use(logger.New())
	app.Use(recover.New())

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	// Init services and handlers
	userSvc := service.NewUserService(database.DB)
	authHandler := handler.NewAuthHandler(userSvc)
	accountHandler := handler.NewAccountHandler()
	meHandler := handler.NewMeHandler(roleClient)

	// ── Auth Routes ───────────────────────────────────────
	// Public
	auth := app.Group("/api/v1/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)
	auth.Post("/login-with-otp", authHandler.LoginWithOTP)
	auth.Post("/refresh-token", authHandler.RefreshToken)

	// Protected
	authProtected := app.Group("/api/v1/auth", middleware.ExtractBearerToken)
	authProtected.Post("/logout", authHandler.Logout)
	authProtected.Put("/update-password", authHandler.UpdatePassword)

	// ── Account Routes ────────────────────────────────────
	// Public
	account := app.Group("/api/v1/account")
	account.Post("/check-exist-value", accountHandler.CheckExistValue)

	// Protected
	accountProtected := app.Group("/api/v1/account", middleware.ExtractBearerToken)
	accountProtected.Get("/information", accountHandler.GetInformation)
	accountProtected.Post("/link-cis", accountHandler.LinkCIS)
	accountProtected.Put("/update-username", accountHandler.UpdateUsername)
	accountProtected.Put("/update-profile", accountHandler.UpdateProfile)
	accountProtected.Post("/get-cis-number", accountHandler.GetCISNumber)

	// ── Me Routes (permissions + menus for frontend) ──────
	// token → accountId resolved via cache/auth-center
	resolveID := middleware.ResolveAccountID(tokenCache)
	me := app.Group("/api/v1/me", middleware.ExtractBearerToken, resolveID)
	me.Get("/permissions", meHandler.GetPermissions)
	me.Get("/menus", meHandler.GetMenus)

	// ── Demo Routes (permission-protected examples) ────────
	demo := app.Group("/api/v1/demo", middleware.ExtractBearerToken, resolveID)
	demo.Get("/setting",
		middleware.RequirePermission(roleClient, "setting", "view"),
		func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"code": "SUCCESS", "data": "setting page data"})
		},
	)
	demo.Post("/setting",
		middleware.RequirePermission(roleClient, "setting", "create"),
		func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"code": "SUCCESS", "data": "setting created"})
		},
	)
	demo.Put("/setting/:id",
		middleware.RequirePermission(roleClient, "setting", "edit"),
		func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"code": "SUCCESS", "data": "setting updated"})
		},
	)
	demo.Delete("/setting/:id",
		middleware.RequirePermission(roleClient, "setting", "delete"),
		func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"code": "SUCCESS", "data": "setting deleted"})
		},
	)

	log.Printf("Server running on port %s", config.Cfg.AppPort)
	log.Fatal(app.Listen(":" + config.Cfg.AppPort))
}

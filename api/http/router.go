package http

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"project-evermos/internal/config"
	authHandler "project-evermos/internal/todo/handler/auth"
	tokoHandler "project-evermos/internal/todo/handler/toko"
	authRepo "project-evermos/internal/todo/repository/auth"
	storeRepo "project-evermos/internal/todo/repository/toko"
	authService "project-evermos/internal/todo/service/auth"
	tokoService "project-evermos/internal/todo/service/auth/toko"
)

// RegisterRoutes registers HTTP routes for the application.
// This keeps the router setup centralized.
func RegisterRoutes(app *fiber.App, gdb *gorm.DB, cfg *config.Config) {
	// Healthcheck endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	// Auth module wiring
	repo := authRepo.NewRepository(gdb)
	storeR := storeRepo.NewRepository(gdb)
	service := authService.NewService(repo, storeR, cfg)
	h := authHandler.NewHandler(service)

	grp := app.Group("/auth")
	grp.Post("/login", h.Login)
	grp.Post("/register", h.Register)

	// Toko module wiring
	tSvc := tokoService.NewService(storeR)
	tH := tokoHandler.NewHandler(tSvc)

	// Protected Toko endpoints (require JWT)
	jwtMW := tokoHandler.JWTMiddleware(cfg.JWTSecret)
	// Register static route before parameterized route to avoid capture as :id_toko
	app.Get("/toko/my", jwtMW, tH.GetMy)

	// Public Toko endpoints
	app.Get("/toko", tH.List)
	app.Get("/toko/:id_toko", tH.GetByID)

	// Protected update endpoint
	app.Put("/toko/:id_toko", jwtMW, tH.Update)
}
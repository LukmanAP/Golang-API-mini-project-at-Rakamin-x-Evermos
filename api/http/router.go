package http

import "github.com/gofiber/fiber/v2"

// RegisterRoutes registers HTTP routes for the application.
// This keeps the router setup centralized.
func RegisterRoutes(app *fiber.App) {
    // Healthcheck endpoint
    app.Get("/health", func(c *fiber.Ctx) error {
        return c.SendString("OK")
    })
}
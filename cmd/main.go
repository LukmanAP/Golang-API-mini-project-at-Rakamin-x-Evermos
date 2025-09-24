package main

import (
    "log"

    "github.com/gofiber/fiber/v2"
    httpRouter "project-evermos/api/http"
    "project-evermos/internal/config"
    "project-evermos/internal/db"
)

func main() {
    // Load config from .env
    cfg, err := config.Load()
    if err != nil {
        log.Fatal(err)
    }

    // Init MySQL connection
    _, err = db.NewMySQL(cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPass, cfg.DBName)
    if err != nil {
        log.Fatal(err)
    }

    // Fiber app
    app := fiber.New()

    // Register centralized routes
    httpRouter.RegisterRoutes(app)

    app.Get("/", func(c *fiber.Ctx) error {
        return c.SendString("hello world")
    })

    if err := app.Listen(":" + cfg.AppPort); err != nil {
        log.Fatal(err)
    }
}

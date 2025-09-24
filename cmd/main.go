package main

import (
    "log"
    "path/filepath"

    "github.com/gofiber/fiber/v2"
    httpRouter "project-evermos/api/http"
    "project-evermos/internal/config"
    "project-evermos/internal/db"
)

func main() {
    cfg, err := config.Load()
    if err != nil {
        log.Fatal(err)
    }

    gdb, err := db.NewMySQL(cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPass, cfg.DBName)
    if err != nil {
        log.Fatal(err)
    }

    // Run migrations from ./migrations
    if err := db.RunMigrations(gdb, filepath.Join(".", "migrations")); err != nil {
        log.Fatal(err)
    }

    app := fiber.New()

    httpRouter.RegisterRoutes(app)

    app.Get("/", func(c *fiber.Ctx) error { return c.SendString("hello world") })

    if err := app.Listen(":" + cfg.AppPort); err != nil {
        log.Fatal(err)
    }
}

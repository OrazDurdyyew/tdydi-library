package main

import (
	"log/slog"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"library/internal/auth"
	"library/internal/config"
	"library/internal/routes"
	"library/internal/services"
	"library/internal/storage"
)

func main() {
	cfg := config.Load()

	db, err := storage.InitDB(cfg.DatabaseURL)
	if err != nil {
		slog.Error("db init failed", "err", err)
		os.Exit(1)
	}
	slog.Info("db connected")

	if err := storage.RunMigrations(db, cfg.MigrationsPath); err != nil {
		slog.Error("migrations failed", "err", err)
		os.Exit(1)
	}
	slog.Info("migrations applied")

	tm := auth.NewTokenManager(cfg.JWTSecret, cfg.JWTExpiry, cfg.RefreshExpiry)
	files := services.NewFileStorage(cfg.StoragePath)

	app := fiber.New(fiber.Config{
		AppName:               "library-api",
		DisableStartupMessage: false,
		BodyLimit:             0,
		StreamRequestBody:     true,
		ReadTimeout:           0,
		WriteTimeout:          0,
	})

	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: cfg.AllowedOrigins,
		AllowMethods: "GET,POST,PUT,PATCH,DELETE,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization",
	}))

	routes.Register(app, db, tm, files, cfg.StoragePath)

	slog.Info("starting server", "port", cfg.APIPort)
	if err := app.Listen(":" + cfg.APIPort); err != nil {
		slog.Error("server failed", "err", err)
		os.Exit(1)
	}
}

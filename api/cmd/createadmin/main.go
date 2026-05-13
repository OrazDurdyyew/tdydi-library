package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"library/internal/auth"
	"library/internal/config"
	"library/internal/models"
	"library/internal/storage"
)

// CLI для создания администратора.
//
// Использование (из контейнера):
//   go run ./cmd/createadmin -username=admin -password=secret -role=superadmin
//
// Или собранным бинарём:
//   /app/library-createadmin -username=admin -password=secret -role=superadmin

func main() {
	username := flag.String("username", "", "username (required)")
	password := flag.String("password", "", "password (required)")
	email := flag.String("email", "", "email (optional)")
	role := flag.String("role", "superadmin", "role: superadmin | content_manager")
	flag.Parse()

	if *username == "" || *password == "" {
		fmt.Fprintln(os.Stderr, "usage: createadmin -username=NAME -password=PASS [-email=EMAIL] [-role=superadmin|content_manager]")
		os.Exit(1)
	}
	if *role != "superadmin" && *role != "content_manager" {
		fmt.Fprintln(os.Stderr, "invalid role: must be 'superadmin' or 'content_manager'")
		os.Exit(1)
	}

	cfg := config.Load()
	db, err := storage.InitDB(cfg.DatabaseURL)
	if err != nil {
		slog.Error("db init failed", "err", err)
		os.Exit(1)
	}

	hash, err := auth.HashPassword(*password)
	if err != nil {
		slog.Error("password hash failed", "err", err)
		os.Exit(1)
	}

	admin := models.Admin{
		Username:     *username,
		PasswordHash: hash,
		Email:        *email,
		Role:         models.AdminRole(*role),
		AllSections:  *role == "superadmin",
	}

	if err := db.Create(&admin).Error; err != nil {
		slog.Error("admin create failed", "err", err)
		os.Exit(1)
	}

	fmt.Printf("admin created: id=%d username=%s role=%s\n", admin.ID, admin.Username, admin.Role)
}

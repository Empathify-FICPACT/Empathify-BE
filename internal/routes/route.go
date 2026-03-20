package routes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Empathify-FICPACT/Empathify-BE/internal/handler"
	"github.com/Empathify-FICPACT/Empathify-BE/internal/provider"
	"github.com/Empathify-FICPACT/Empathify-BE/internal/repository"
	"github.com/Empathify-FICPACT/Empathify-BE/internal/service"
)

func Route(app *fiber.App, db *pgxpool.Pool) {
	authRepo := repository.NewAuthRepository(db)
	authService := service.NewAuthService(authRepo)
	googleOAuth := provider.NewGoogleOAuthProvider()
	authHandler := handler.NewAuthHandler(authService, googleOAuth)

	api := app.Group("/api/v1")

	api.Get("/", func (c fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	RegisterAuthRoutes(api, authHandler)
}
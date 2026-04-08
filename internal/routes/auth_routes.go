package routes

import (
	"github.com/gofiber/fiber/v3"

	"github.com/Empathify-FICPACT/Empathify-BE/internal/handler"
)

func RegisterAuthRoutes(router fiber.Router, authHandler *handler.AuthHandler) {
	auth := router.Group("/auth")

	// email
	auth.Post("/register", authHandler.RegisterEmail)
	auth.Post("/login", authHandler.LoginEmail)

	// google oauth
	auth.Get("/google", authHandler.GoogleLogin)
	auth.Get("/google/callback", authHandler.GoogleCallback)
	auth.Get("/me", authHandler.Me)
}
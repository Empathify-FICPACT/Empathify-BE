package routes

import (
	"github.com/gofiber/fiber/v3"

	"github.com/Empathify-FICPACT/Empathify-BE/internal/handler"
	"github.com/Empathify-FICPACT/Empathify-BE/internal/middleware"
)

func RegisterUserRoutes(router fiber.Router, userHandler *handler.UserHandler) {
	user := router.Group("/user", middleware.AuthMiddleware())

	user.Patch("/onboarding", userHandler.Onboarding)
}
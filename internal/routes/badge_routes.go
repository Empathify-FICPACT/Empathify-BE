package routes

import (
	"github.com/gofiber/fiber/v3"

	"github.com/Empathify-FICPACT/Empathify-BE/internal/handler"
	"github.com/Empathify-FICPACT/Empathify-BE/internal/middleware"
)

func RegisterBadgeRoutes(router fiber.Router, badgeHandler *handler.BadgeHandler) {
	badge := router.Group("/badges", middleware.AuthMiddleware())

	badge.Get("", badgeHandler.GetUserBadges)
}
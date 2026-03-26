package routes

import (
	"github.com/gofiber/fiber/v3"

	"github.com/Empathify-FICPACT/Empathify-BE/internal/handler"
	"github.com/Empathify-FICPACT/Empathify-BE/internal/middleware"
)

func RegisterDashboardRoutes(router fiber.Router, dashboardHandler *handler.DashboardHandler) {
	router.Get("/dashboard", middleware.AuthMiddleware(), dashboardHandler.GetDashboard)
}
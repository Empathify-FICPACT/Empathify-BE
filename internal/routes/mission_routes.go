package routes

import (
	"github.com/gofiber/fiber/v3"

	"github.com/Empathify-FICPACT/Empathify-BE/internal/handler"
	"github.com/Empathify-FICPACT/Empathify-BE/internal/middleware"
)

func RegisterMissionRoutes(router fiber.Router, missionHandler *handler.MissionHandler) {
	mission := router.Group("/missions", middleware.AuthMiddleware())

	mission.Get("/today", missionHandler.GetTodayMissions)
}
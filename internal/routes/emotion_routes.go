package routes

import (
	"github.com/gofiber/fiber/v3"

	"github.com/Empathify-FICPACT/Empathify-BE/internal/handler"
	"github.com/Empathify-FICPACT/Empathify-BE/internal/middleware"
)

func RegisterEmotionRoutes(router fiber.Router, emotionHandler *handler.EmotionHandler) {
	emotion := router.Group("/emotion", middleware.AuthMiddleware())

	emotion.Post("/sessions", emotionHandler.StartSession)
	emotion.Post("/sessions/:id/answers", emotionHandler.SubmitAnswer)
	emotion.Patch("/sessions/:id/complete", emotionHandler.CompleteSession)
}
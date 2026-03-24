package routes

import (
	"github.com/gofiber/fiber/v3"

	"github.com/Empathify-FICPACT/Empathify-BE/internal/handler"
	"github.com/Empathify-FICPACT/Empathify-BE/internal/middleware"
)

func RegisterStoryRoutes(router fiber.Router, storyHandler *handler.StoryHandler) {
	story := router.Group("/story", middleware.AuthMiddleware())

	story.Post("/sessions", storyHandler.StartSession)
	story.Post("/sessions/:id/answers", storyHandler.SubmitAnswer)
	story.Patch("/sessions/:id/complete", storyHandler.CompleteSession)
}
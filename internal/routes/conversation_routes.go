package routes

import (
	"github.com/gofiber/fiber/v3"

	"github.com/Empathify-FICPACT/Empathify-BE/internal/handler"
	"github.com/Empathify-FICPACT/Empathify-BE/internal/middleware"
)

func RegisterConversationRoutes(router fiber.Router, convHandler *handler.ConversationHandler) {
	conv := router.Group("/conversation", middleware.AuthMiddleware())

	conv.Get("/topics", convHandler.GetTopics)
	conv.Post("/sessions", convHandler.StartSession)
	conv.Post("/sessions/:id/messages", convHandler.SendMessage)
	conv.Patch("/sessions/:id/complete", convHandler.CompleteSession)
}
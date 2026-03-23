package routes

import (
	"github.com/gofiber/fiber/v3"

	"github.com/Empathify-FICPACT/Empathify-BE/internal/handler"
	"github.com/Empathify-FICPACT/Empathify-BE/internal/middleware"
)

func RegisterExpressionRoutes(router fiber.Router, exprHandler *handler.ExpressionHandler) {
	expr := router.Group("/expression", middleware.AuthMiddleware())

	expr.Post("/sessions", exprHandler.StartSession)
	expr.Post("/sessions/:id/attempts", exprHandler.SubmitAttempt)
	expr.Patch("/sessions/:id/complete", exprHandler.CompleteSession)
}
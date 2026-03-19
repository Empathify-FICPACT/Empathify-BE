package routes

import "github.com/gofiber/fiber/v3"

func Route(app *fiber.App) {
	api := app.Group("/api/v1")

	api.Get("/", func (c fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})
}
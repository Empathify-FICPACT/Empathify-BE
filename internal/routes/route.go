package routes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Route(app *fiber.App, db *pgxpool.Pool) {
	api := app.Group("/api/v1")

	api.Get("/", func (c fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})
}
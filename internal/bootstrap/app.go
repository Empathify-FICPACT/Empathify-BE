package bootstrap

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Empathify-FICPACT/Empathify-BE/internal/routes"
)

func NewApp(db *pgxpool.Pool) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName: "Empathify API",
	})

	app.Use(logger.New())

	app.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:3000",
			"https://empathifyy.vercel.app",
		},
		AllowMethods: []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))


	routes.Route(app, db)

	return app
}
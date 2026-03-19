package bootstrap

import "github.com/gofiber/fiber/v3"

func InitApp() *fiber.App {
	app := fiber.New()

	return app
}
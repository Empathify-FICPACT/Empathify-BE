package routes

import (
	swaggo "github.com/gofiber/contrib/v3/swaggo"
	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Empathify-FICPACT/Empathify-BE/internal/handler"
	"github.com/Empathify-FICPACT/Empathify-BE/internal/provider"
	"github.com/Empathify-FICPACT/Empathify-BE/internal/repository"
	"github.com/Empathify-FICPACT/Empathify-BE/internal/service"
)

func Route(app *fiber.App, db *pgxpool.Pool) {
	authRepo := repository.NewAuthRepository(db)
	authService := service.NewAuthService(authRepo)
	googleOAuth := provider.NewGoogleOAuthProvider()
	authHandler := handler.NewAuthHandler(authService, googleOAuth)

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	stt := provider.NewSTTProvider()
	gemini := provider.NewGeminiProvider()
	convRepo := repository.NewConversationRepository(db)
	convService := service.NewConversationService(convRepo, userRepo, gemini, stt)
	convHandler := handler.NewConversationHandler(convService)

	exprRepo := repository.NewExpressionRepository(db)
	exprService := service.NewExpressionService(exprRepo, userRepo, gemini)
	exprHandler := handler.NewExpressionHandler(exprService)

	emotionRepo := repository.NewEmotionRepository(db)
	emotionService := service.NewEmotionService(emotionRepo, userRepo)
	emotionHandler := handler.NewEmotionHandler(emotionService)

	storyRepo := repository.NewStoryRepository(db)
	storyService := service.NewStoryService(storyRepo, userRepo)
	storyHandler := handler.NewStoryHandler(storyService)

	api := app.Group("/api/v1")

	api.Get("/", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	app.Get("/swagger/*", swaggo.HandlerDefault)

	RegisterAuthRoutes(api, authHandler)
	RegisterUserRoutes(api, userHandler)
	RegisterConversationRoutes(api, convHandler)
	RegisterExpressionRoutes(api, exprHandler)
	RegisterEmotionRoutes(api, emotionHandler)
	RegisterStoryRoutes(api, storyHandler)
}
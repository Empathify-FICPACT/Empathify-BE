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
	// repos
	authRepo    := repository.NewAuthRepository(db)
	userRepo    := repository.NewUserRepository(db)
	convRepo    := repository.NewConversationRepository(db)
	exprRepo    := repository.NewExpressionRepository(db)
	emotionRepo := repository.NewEmotionRepository(db)
	storyRepo   := repository.NewStoryRepository(db)
	missionRepo := repository.NewMissionRepository(db)
	badgeRepo := repository.NewBadgeRepository(db)
	dashboardRepo    := repository.NewDashboardRepository(db)

	// providers
	gemini      := provider.NewGeminiProvider()
	stt         := provider.NewSTTProvider()
	googleOAuth := provider.NewGoogleOAuthProvider()

	// mission service dibuat duluan karena dipakai service lain
	missionSvc := service.NewMissionService(missionRepo, userRepo)
	badgeSvc := service.NewBadgeService(badgeRepo, userRepo, exprRepo, emotionRepo, storyRepo, convRepo)

	// services
	authSvc    := service.NewAuthService(authRepo)
	userSvc    := service.NewUserService(userRepo)
	convSvc    := service.NewConversationService(convRepo, userRepo, gemini, stt, missionSvc, badgeSvc)
	exprSvc    := service.NewExpressionService(exprRepo, userRepo, gemini, missionSvc, badgeSvc)
	emotionSvc := service.NewEmotionService(emotionRepo, userRepo, missionSvc, badgeSvc)
	storySvc   := service.NewStoryService(storyRepo, userRepo, missionSvc, badgeSvc)
	dashboardSvc     := service.NewDashboardService(dashboardRepo, userRepo)

	// handlers
	authHandler    := handler.NewAuthHandler(authSvc, googleOAuth)
	userHandler    := handler.NewUserHandler(userSvc)
	convHandler    := handler.NewConversationHandler(convSvc)
	exprHandler    := handler.NewExpressionHandler(exprSvc)
	emotionHandler := handler.NewEmotionHandler(emotionSvc)
	storyHandler   := handler.NewStoryHandler(storySvc)
	missionHandler := handler.NewMissionHandler(missionSvc)
	badgeHandler := handler.NewBadgeHandler(badgeSvc)
	dashboardHandler := handler.NewDashboardHandler(dashboardSvc)

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
	RegisterMissionRoutes(api, missionHandler)
	RegisterBadgeRoutes(api, badgeHandler)
	RegisterDashboardRoutes(api, dashboardHandler)
}
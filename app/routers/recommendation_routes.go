package routes

import (
	"SyncBeat/config"
	"SyncBeat/controllers"
	"SyncBeat/middleware"
	"SyncBeat/repositories"
	"SyncBeat/service"

	"github.com/gin-gonic/gin"
)

func SetupRecommendationRoutes(router *gin.Engine, cfg *config.Config) {
	recommendationRepo := repositories.NewRecommendationRepository()
	recommendationService := service.NewRecommendationService(recommendationRepo)
	recommendationController := controllers.NewRecommendationController(recommendationService)

	// ミドルウェア
	authMiddleware, err := middleware.AuthMiddleware(cfg)
	if err != nil {
		panic(err)
	}

	recommendations := router.Group("/app/recommendations")
	{
		// 認証が必要なルート
		recommendations.Use(authMiddleware.MiddlewareFunc())
		{
			recommendations.POST("/", recommendationController.GetRecommendation)
			recommendations.PUT("/feedback", recommendationController.PutRecommendationFeedback)
		}
	}
}

package recommandation

import (
	"SyncBeat/config"
	recommendation2 "SyncBeat/controllers/recommendation"
	"SyncBeat/middleware/jwt"
	"SyncBeat/repositories/recommendation"
	"SyncBeat/service/recommandation"

	"github.com/gin-gonic/gin"
)

func SetupRecommendationRoutes(router *gin.Engine, cfg *config.Config) {
	recommendationRepo := recommendation.NewRecommendationRepository()
	recommendationService := recommandation.NewRecommendationService(recommendationRepo)
	recommendationController := recommendation2.NewRecommendationController(recommendationService)

	// ミドルウェア
	authMiddleware, err := jwt.AuthMiddleware(cfg)
	if err != nil {
		panic(err)
	}

	recommendations := router.Group("/api/recommendations")
	{
		// 認証が必要なルート
		recommendations.Use(authMiddleware.MiddlewareFunc())
		{
			recommendations.POST("/", recommendationController.GetRecommendation)
			recommendations.PUT("/feedback", recommendationController.PutRecommendationFeedback)
		}
	}
}

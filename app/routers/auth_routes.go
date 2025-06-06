package routes

import (
	"SyncBeat/config"
	"SyncBeat/controllers"
	"SyncBeat/middleware"
	"github.com/gin-gonic/gin"
)

func SetupAuthRoutes(router *gin.Engine, cfg *config.Config) {
	authController := controllers.NewAuthController(cfg)

	// ミドルウェア
	authMiddleware, err := middleware.AuthMiddleware(cfg)
	if err != nil {
		panic(err)
	}

	auth := router.Group("/api/auth")
	{
		auth.POST("/register", authController.Register)
		auth.POST("/login", authMiddleware.LoginHandler)

		// 認証が必要なルート
		auth.Use(authMiddleware.MiddlewareFunc())
		{
			auth.GET("/refresh_token", authMiddleware.RefreshHandler)
			auth.GET("/logout", authMiddleware.LogoutHandler)
		}
	}
}

package auth

import (
	"SyncBeat/config"
	"SyncBeat/controllers/auth"
	"SyncBeat/middleware/jwt"
	"github.com/gin-gonic/gin"
)

func SetupAuthRoutes(router *gin.Engine, cfg *config.Config) {
	authController := auth.NewAuthController(cfg)

	// ミドルウェア
	authMiddleware, err := jwt.AuthMiddleware(cfg)
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

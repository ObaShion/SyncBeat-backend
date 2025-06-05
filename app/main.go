package main

import (
	database "SyncBeat/migrations"
	routes "SyncBeat/routers"
	"github.com/joho/godotenv"
	"log"
	"os"

	"SyncBeat/config"

	"github.com/gin-gonic/gin"
)

func main() {
	// .envファイルの読み込み
	if err := godotenv.Load(); err != nil {
		log.Println(".env not found")
	}

	// 設定の読み込み
	cfg := config.LoadConfig()

	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// データベース初期化
	database.InitDB()

	// ルーター設定
	router := gin.Default()

	// ルート設定
	routes.SetupAuthRoutes(router, cfg)
	routes.SetupRecommendationRoutes(router, cfg)

	if err := router.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

package database

import (
	"SyncBeat/models/auth"
	"SyncBeat/models/recommendation"
	"SyncBeat/models/user"
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB() {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		log.Fatalf("データベース接続に失敗しました: %v", err)
	}

	// マイグレーション
	err = DB.AutoMigrate(
		&user.User{},
		&user.UserState{},
		&recommendation.MusicRecommendation{},
		&auth.RefreshToken{},
	)
	if err != nil {
		log.Fatalf("マイグレーションに失敗しました: %v", err)
	}

	log.Println("データベース接続が成功しました")
}

package auth

import (
	"SyncBeat/config"
	"SyncBeat/errors"
	database "SyncBeat/migrations"
	"SyncBeat/models/user"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthController struct {
	config *config.Config
}

func NewAuthController(cfg *config.Config) *AuthController {
	return &AuthController{
		config: cfg,
	}
}

type RegisterInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

func (ac *AuthController) Register(c *gin.Context) {
	var input RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(errors.NewBadRequest("無効なリクエストです", err).Code, gin.H{
			"error": err.Error(),
		})
		return
	}

	user := user.User{
		Email:    input.Email,
		Password: input.Password,
		UID:      uuid.NewString(),
	}

	if err := user.HashPassword(); err != nil {
		c.JSON(errors.NewInternalServerError("パスワードのハッシュ化に失敗しました", err).Code, gin.H{
			"error": "パスワードのハッシュ化に失敗しました",
		})
		return
	}

	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(errors.NewBadRequest("ユーザーの作成に失敗しました", err).Code, gin.H{
			"error": "ユーザーの作成に失敗しました",
		})
		return
	}

	c.JSON(201, gin.H{
		"message": "ユーザーが正常に作成されました",
	})
}

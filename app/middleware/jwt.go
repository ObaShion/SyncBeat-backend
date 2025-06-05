package middleware

import (
	"SyncBeat/config"
	"SyncBeat/errors"
	database "SyncBeat/migrations"
	models "SyncBeat/models"
	"net/http"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type login struct {
	Email    string `form:"email" json:"email" binding:"required,email"`
	Password string `form:"password" json:"password" binding:"required"`
}

type RefreshToken struct {
	gorm.Model
	UserID    uint      `gorm:"not null"`
	Token     string    `gorm:"uniqueIndex;not null"`
	ExpiresAt time.Time `gorm:"not null"`
}

func AuthMiddleware(cfg *config.Config) (*jwt.GinJWTMiddleware, error) {
	return jwt.New(&jwt.GinJWTMiddleware{
		Realm:       cfg.JWT.Realm,
		Key:         []byte(cfg.JWT.SecretKey),
		Timeout:     cfg.JWT.Timeout,    // 24時間
		MaxRefresh:  cfg.JWT.MaxRefresh, // 7日間
		IdentityKey: cfg.JWT.IdentityKey,

		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*models.User); ok {
				return jwt.MapClaims{
					"id":    v.ID,
					"email": v.Email,
				}
			}
			return jwt.MapClaims{}
		},

		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &models.User{
				Model: gorm.Model{ID: uint(claims["id"].(float64))},
				Email: claims["email"].(string),
			}
		},

		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginVals login
			if err := c.ShouldBind(&loginVals); err != nil {
				return nil, errors.NewBadRequest("無効なリクエストです", err)
			}

			var user models.User
			if err := database.DB.Where("email = ?", loginVals.Email).First(&user).Error; err != nil {
				return nil, errors.NewUnauthorized("メールアドレスまたはパスワードが正しくありません", err)
			}

			if err := user.CheckPassword(loginVals.Password); err != nil {
				return nil, errors.NewUnauthorized("メールアドレスまたはパスワードが正しくありません", err)
			}

			return &user, nil
		},

		Authorizator: func(data interface{}, c *gin.Context) bool {
			if _, ok := data.(*models.User); ok {
				return true
			}
			return false
		},

		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},

		TokenLookup:   "header: Authorization, query: token, cookie: jwt",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,

		SendCookie:     true,
		SecureCookie:   true,
		CookieHTTPOnly: true,
		CookieName:     "jwt",
		CookieSameSite: http.SameSiteStrictMode,

		// リフレッシュトークンの処理
		RefreshResponse: func(c *gin.Context, code int, token string, expire time.Time) {
			// リフレッシュトークンをデータベースに保存
			claims := jwt.ExtractClaims(c)
			userID := uint(claims["id"].(float64))

			refreshToken := RefreshToken{
				UserID:    userID,
				Token:     token,
				ExpiresAt: expire,
			}

			if err := database.DB.Create(&refreshToken).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"code":    http.StatusInternalServerError,
					"message": "リフレッシュトークンの保存に失敗しました",
				})
				return
			}

			c.JSON(code, gin.H{
				"code":          code,
				"token":         token,
				"expire":        expire.Format(time.RFC3339),
				"token_type":    "Bearer",
				"expires_in":    int(expire.Sub(time.Now()).Seconds()),
				"refresh_token": refreshToken.Token,
			})
		},
	})
}

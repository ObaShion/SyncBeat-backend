package controllers

import (
	"SyncBeat/errors"
	"SyncBeat/service"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

type RecommendationController struct {
	service *service.RecommendationService
}

func NewRecommendationController(service *service.RecommendationService) *RecommendationController {
	return &RecommendationController{
		service: service,
	}
}

func (rc *RecommendationController) GetRecommendation(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	userID, ok := claims["id"].(float64)
	if !ok {
		c.JSON(errors.NewUnauthorized("無効なユーザー情報です", nil).Code, gin.H{
			"error": "無効なユーザー情報です",
		})
		return
	}

	var input service.UserStateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(errors.NewBadRequest("無効なリクエストです", err).Code, gin.H{
			"error": err.Error(),
		})
		return
	}

	response, err := rc.service.GetRecommendation(uint(userID), input)
	if err != nil {
		c.JSON(errors.NewInternalServerError("推薦の生成に失敗しました", err).Code, gin.H{
			"error": "推薦の生成に失敗しました",
		})
		return
	}

	c.JSON(200, response)
}

func (rc *RecommendationController) PutRecommendationFeedback(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	userID, ok := claims["id"].(float64)
	if !ok {
		c.JSON(errors.NewUnauthorized("無効なユーザー情報です", nil).Code, gin.H{
			"error": "無効なユーザー情報です",
		})
		return
	}

	var input struct {
		UID   string `json:"uid" binding:"required"`
		Score int    `json:"score" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(errors.NewBadRequest("無効なリクエストです", err).Code, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := rc.service.UpdateRecommendationScore(input.UID, uint(userID), input.Score); err != nil {
		c.JSON(errors.NewInternalServerError("推薦の更新に失敗しました", err).Code, gin.H{
			"error": "推薦の更新に失敗しました",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "推薦のスコアを更新しました",
	})
}

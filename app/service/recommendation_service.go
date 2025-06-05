package service

import (
	"SyncBeat/models"
	"SyncBeat/repositories"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/google/uuid"
)

type RecommendationService struct {
	repo *repositories.RecommendationRepository
}

func NewRecommendationService(repo *repositories.RecommendationRepository) *RecommendationService {
	return &RecommendationService{
		repo: repo,
	}
}

type UserStateInput struct {
	HeartRates []int    `json:"heart_rates" binding:"required"`
	Weather    string   `json:"weather" binding:"required"`
	Movement   string   `json:"movement" binding:"required"`
	Latitude   float64  `json:"latitude" binding:"required"`
	Longitude  float64  `json:"longitude" binding:"required"`
	Calendar   []string `json:"calendar"`
	MusicLimit int      `json:"music_limit"`
}

type RecommendationItem struct {
	MusicID string `json:"music_id"`
	UID     string `json:"uid"`
}

type MusicRecommendationResponse struct {
	Recommendations []RecommendationItem `json:"recommendations"`
}

func (rs *RecommendationService) GetRecommendation(userID uint, input UserStateInput) (*MusicRecommendationResponse, error) {
	musicLimit := input.MusicLimit
	if musicLimit <= 0 {
		musicLimit = 2 // デフォルト値
	}

	heartRates := make([]int64, len(input.HeartRates))
	for i, v := range input.HeartRates {
		heartRates[i] = int64(v)
	}

	userState, err := rs.repo.CreateUserState(userID, heartRates, input.Weather, input.Movement, input.Latitude, input.Longitude, input.Calendar)
	if err != nil {
		return nil, err
	}

	recommendationItems := rs.generateRecommendation(*userState, musicLimit)

	// 推薦を保存
	for i := range recommendationItems {
		recommendation, err := rs.repo.CreateMusicRecommendation(userID, userState.ID, recommendationItems[i].MusicID)
		if err != nil {
			return nil, err
		}
		recommendationItems[i].UID = recommendation.UID
	}

	return &MusicRecommendationResponse{
		Recommendations: recommendationItems,
	}, nil
}

func (rs *RecommendationService) UpdateRecommendationScore(uid string, userID uint, score int) error {
	return rs.repo.UpdateRecommendationScore(uid, userID, score)
}

// generateRecommendation
func (rs *RecommendationService) generateRecommendation(state models.UserState, musicLimit int) []RecommendationItem {

	// DifyAPIのエンドポイント
	url := os.Getenv("DIFY_BASE_URL")

	// 心拍数を文字列に変換
	heartRatesStr := ""
	for i, rate := range state.HeartRates {
		if i > 0 {
			heartRatesStr += ","
		}
		heartRatesStr += fmt.Sprintf("%d", rate)
	}

	// カレンダーを文字列に変換
	calendarStr := ""
	for i, event := range state.Calendar {
		if i > 0 {
			calendarStr += ","
		}
		calendarStr += event
	}

	// リクエストボディの作成
	requestBody := map[string]interface{}{
		"inputs": map[string]interface{}{
			"heart_rates": heartRatesStr,
			"weather":     state.Weather,
			"movement":    state.Movement,
			"calendar":    calendarStr,
			"music_limit": musicLimit,
		},
		"response_mode": "blocking",
		"user":          "user_1",
	}

	// JSONに変換
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return []RecommendationItem{{MusicID: "error", UID: "error"}}
	}

	// HTTPリクエストの作成
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return []RecommendationItem{{MusicID: "error", UID: "error"}}
	}

	// ヘッダーの設定
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+os.Getenv("DIFY_API_KEY"))

	// リクエストの送信
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return []RecommendationItem{{MusicID: "error", UID: "error"}}
	}

	defer resp.Body.Close()

	// レスポンスの読み取り
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return []RecommendationItem{{MusicID: "error", UID: "error"}}
	}

	// レスポンスからデータを取得
	data, ok := result["data"].(map[string]interface{})
	if !ok {
		return []RecommendationItem{{MusicID: "error", UID: "error"}}
	}

	// outputsから音楽ID配列を取得
	outputs, ok := data["outputs"].(map[string]interface{})
	if !ok {
		return []RecommendationItem{{MusicID: "error", UID: "error"}}
	}

	musicIDsInterface, ok := outputs["music_ids"].([]interface{})
	if !ok {
		return []RecommendationItem{{MusicID: "error", UID: "error"}}
	}

	// インターフェース配列を文字列配列に変換
	musicIDs := make([]string, len(musicIDsInterface))
	for i, v := range musicIDsInterface {
		if str, ok := v.(string); ok {
			musicIDs[i] = str
		} else {
			return []RecommendationItem{{MusicID: "error", UID: "error"}}
		}
	}

	recommendationItems := make([]RecommendationItem, 0, len(musicIDs))
	for _, musicID := range musicIDs {
		recommendationItems = append(recommendationItems, RecommendationItem{
			MusicID: musicID,
			UID:     uuid.NewString(),
		})
	}
	return recommendationItems
}

package repositories

import (
	database "SyncBeat/migrations"
	"SyncBeat/models"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type RecommendationRepository struct{}

func NewRecommendationRepository() *RecommendationRepository {
	return &RecommendationRepository{}
}

func (rr *RecommendationRepository) CreateUserState(userID uint, heartRates []int64, weather, movement string, latitude, longitude float64, calendar []string) (*models.UserState, error) {
	userState := models.UserState{
		UserID:     userID,
		HeartRates: pq.Int64Array(heartRates),
		Weather:    weather,
		Movement:   movement,
		Latitude:   latitude,
		Longitude:  longitude,
		Calendar:   pq.StringArray(calendar),
	}

	if err := database.DB.Create(&userState).Error; err != nil {
		return nil, err
	}
	return &userState, nil
}

func (rr *RecommendationRepository) CreateMusicRecommendation(userID uint, userStateID uint, musicID string) (*models.MusicRecommendation, error) {
	uid := uuid.NewString()
	musicRecommendation := models.MusicRecommendation{
		UserID:      userID,
		UserStateID: userStateID,
		MusicID:     musicID,
		UID:         uid,
	}

	if err := database.DB.Create(&musicRecommendation).Error; err != nil {
		return nil, err
	}
	return &musicRecommendation, nil
}

func (rr *RecommendationRepository) UpdateRecommendationScore(uid string, userID uint, score int) error {
	var recommendation models.MusicRecommendation
	if err := database.DB.Where("uid = ? AND user_id = ?", uid, userID).First(&recommendation).Error; err != nil {
		return err
	}

	recommendation.Score = score
	return database.DB.Save(&recommendation).Error
}

package models

import (
	"gorm.io/gorm"
)

// MusicRecommendation 音楽の推薦を表すモデル
type MusicRecommendation struct {
	gorm.Model
	UserID      uint   `gorm:"not null" json:"user_id"`
	UserStateID uint   `gorm:"not null" json:"user_state_id"`
	UID         string `gorm:"uniqueIndex;not null" json:"uid"`
	MusicID     string `json:"music_id"`
	Score       int    `json:"score"`
}

package user

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

// UserState ユーザーの状態を表すモデル
type UserState struct {
	gorm.Model
	UserID     uint           `gorm:"not null" json:"user_id"`
	UID        string         `gorm:"uniqueIndex;not null" json:"uid"`
	HeartRates pq.Int64Array  `gorm:"type:integer[]" json:"heart_rates"` // 過去1時間の心拍数
	Weather    string         `json:"weather"`                           // 天気
	Movement   string         `json:"movement"`                          // 移動状態
	Latitude   float64        `json:"latitude"`                          // 緯度
	Longitude  float64        `json:"longitude"`                         // 経度
	Calendar   pq.StringArray `gorm:"type:text[]" json:"calendar"`       // その日のカレンダー
	Prompt     string         `json:"prompt"`                            // プロンプトがある場合
}

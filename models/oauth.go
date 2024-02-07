package models

import (
	"gorm.io/gorm"
	"time"
)

type OAuth struct {
	gorm.Model
	Provider      	string `gorm:"not null"`               // OAuth 제공자 (예: "Google", "Facebook")
	Provider_userid string `gorm:"unique;not null"`        // OAuth 제공자에서의 사용자 고유 ID
	Access_token 	string `gorm:"not null"`               // 액세스 토큰
	Refresh_token 	string                                 // 리프레시 토큰 (필요에 따라 사용)
	Expiry        	time.Time                              // 액세스 토큰의 만료 시간

	User_id        	uint   `gorm:"uniqueIndex"`                  // Users 테이블의 ID와 연결
}
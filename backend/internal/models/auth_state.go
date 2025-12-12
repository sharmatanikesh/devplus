package models

import "time"

type AuthState struct {
	ID        string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	State     string    `gorm:"uniqueIndex;not null" json:"state"`
	ExpiresAt time.Time `gorm:"not null" json:"expires_at"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func (AuthState) TableName() string {
	return "auth_states"
}

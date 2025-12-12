package models

import "time"

type Session struct {
	ID        string    `gorm:"primaryKey" json:"id"` // UUID
	UserID    int64     `gorm:"index;not null" json:"user_id"`
	User      User      `gorm:"foreignKey:UserID" json:"-"`
	ExpiresAt time.Time `gorm:"not null" json:"expires_at"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func (Session) TableName() string {
	return "sessions"
}

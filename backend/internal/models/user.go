package models

import (
	"time"
)

type User struct {
	ID           string     `gorm:"column:id;primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	CreatedAt    time.Time  `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time  `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	DeletedAt    *time.Time `gorm:"column:deleted_at;index:idx_users_deleted_at" json:"deleted_at"`
	GithubID     int64      `gorm:"column:github_id;uniqueIndex:idx_users_github_id;not null" json:"github_id"`
	Username     string     `gorm:"column:username;not null" json:"username"`
	Email        string     `gorm:"column:email" json:"email"`
	AvatarURL    string     `gorm:"column:avatar_url" json:"avatar_url"`
	AccessToken  string     `gorm:"column:access_token" json:"-"`  // Don't expose in JSON
	RefreshToken string     `gorm:"column:refresh_token" json:"-"` // Don't expose in JSON
}

func (User) TableName() string {
	return "public.users"
}

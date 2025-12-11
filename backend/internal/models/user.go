package models

import (
	"time"
)

type User struct {
	ID           int64      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	CreatedAt    *time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt    *time.Time `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt    *time.Time `gorm:"column:deleted_at;index:idx_users_deleted_at" json:"deleted_at"`
	GithubID     *int64     `gorm:"column:github_id;uniqueIndex:idx_users_github_id" json:"github_id"`
	Username     *string    `gorm:"column:username" json:"username"`
	Email        *string    `gorm:"column:email" json:"email"`
	AccessToken  *string    `gorm:"column:access_token" json:"access_token"`
	RefreshToken *string    `gorm:"column:refresh_token" json:"refresh_token"`
}

func (User) TableName() string {
	return "public.users"
}

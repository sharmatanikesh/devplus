package models

import (
	"time"
)

type User struct {
	ID           int64      `gorm:"column:id;primaryKey;autoIncrement"`
	CreatedAt    *time.Time `gorm:"column:created_at"`
	UpdatedAt    *time.Time `gorm:"column:updated_at"`
	DeletedAt    *time.Time `gorm:"column:deleted_at;index:idx_users_deleted_at"`
	GithubID     *int64     `gorm:"column:github_id;uniqueIndex:idx_users_github_id"`
	Username     *string    `gorm:"column:username"`
	Email        *string    `gorm:"column:email"`
	AccessToken  *string    `gorm:"column:access_token"`
	RefreshToken *string    `gorm:"column:refresh_token"`
}

func (User) TableName() string {
	return "public.users"
}

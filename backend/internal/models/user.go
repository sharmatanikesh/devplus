package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	GithubID     int64  `json:"github_id" gorm:"uniqueIndex"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	AccessToken  string `json:"-"`
	RefreshToken string `json:"-"`
}

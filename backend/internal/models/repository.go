package models

import "gorm.io/gorm"

type Repository struct {
	gorm.Model
	GithubRepoID   int64  `json:"github_repo_id" gorm:"uniqueIndex"`
	Name           string `json:"name"`
	Owner          string `json:"owner"`
	URL            string `json:"url"`
	InstallationID int64  `json:"installation_id"`
}

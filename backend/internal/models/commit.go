package models

import "gorm.io/gorm"

type Commit struct {
	gorm.Model
	SHA        string     `json:"sha" gorm:"uniqueIndex"`
	Message    string     `json:"message"`
	AuthorName string     `json:"author_name"`
	RepoID     uint       `json:"repo_id"`
	Repository Repository `json:"repository" gorm:"foreignKey:RepoID"`
}

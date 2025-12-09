package models

import (
	"time"

	"gorm.io/gorm"
)

type MetricType string

const (
	MetricLeadTime   MetricType = "LEAD_TIME"
	MetricMergeRate  MetricType = "MERGE_RATE"
	MetricReviewTime MetricType = "REVIEW_TIME"
)

type Metric struct {
	gorm.Model
	RepoID     uint       `json:"repo_id"`
	Repository Repository `json:"repository" gorm:"foreignKey:RepoID"`
	Type       MetricType `json:"type"`
	Value      float64    `json:"value"`
	Date       time.Time  `json:"date"`
}

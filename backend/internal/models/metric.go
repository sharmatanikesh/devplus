package models

import (
	"time"
)

type Metric struct {
	ID         int64       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	CreatedAt  *time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt  *time.Time  `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt  *time.Time  `gorm:"column:deleted_at;index:idx_metrics_deleted_at" json:"deleted_at"`
	RepoID     *int64      `gorm:"column:repo_id" json:"repo_id"`
	Repository *Repository `gorm:"foreignKey:RepoID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"repository"`
	Type       *string     `gorm:"column:type" json:"type"`
	Value      *float64    `gorm:"column:value" json:"value"`
	Date       *time.Time  `gorm:"column:date" json:"date"`
}

func (Metric) TableName() string {
	return "public.metrics"
}

package models

import (
	"time"
)

type Metric struct {
	ID         int64       `gorm:"column:id;primaryKey;autoIncrement"`
	CreatedAt  *time.Time  `gorm:"column:created_at"`
	UpdatedAt  *time.Time  `gorm:"column:updated_at"`
	DeletedAt  *time.Time  `gorm:"column:deleted_at;index:idx_metrics_deleted_at"`
	RepoID     *int64      `gorm:"column:repo_id"`
	Repository *Repository `gorm:"foreignKey:RepoID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Type       *string     `gorm:"column:type"`
	Value      *float64    `gorm:"column:value"` // numeric in SQL maps to float64 often, or specialized type
	Date       *time.Time  `gorm:"column:date"`
}

func (Metric) TableName() string {
	return "public.metrics"
}

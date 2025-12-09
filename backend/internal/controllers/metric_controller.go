package controllers

import (
	"net/http"
	"time"

	"devplus-backend/internal/models"

	"github.com/gin-gonic/gin"
)

type MetricController struct{}

func NewMetricController() *MetricController {
	return &MetricController{}
}

func (ctrl *MetricController) GetMetrics(c *gin.Context) {
	// Dummy Metrics
	metrics := []models.Metric{
		{
			Type:  models.MetricLeadTime,
			Value: 12.5, // hours
			Date:  time.Now(),
		},
		{
			Type:  models.MetricMergeRate,
			Value: 0.85, // 85%
			Date:  time.Now(),
		},
	}
	c.JSON(http.StatusOK, metrics)
}

package router

import (
	"devplus-backend/internal/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	// Initialize Controllers
	webhookCtrl := controllers.NewWebhookController()
	repoCtrl := controllers.NewRepoController()
	actionCtrl := controllers.NewActionController()
	metricCtrl := controllers.NewMetricController()

	api := r.Group("/api/v1")
	{
		// Webhooks
		api.POST("/webhook/github", webhookCtrl.HandleWebhook)

		// Repos
		api.GET("/repos", repoCtrl.ListRepos)
		api.POST("/repos/:id/sync", repoCtrl.SyncRepo)
		api.GET("/repos/:id/prs", repoCtrl.ListPRs)
		api.GET("/repos/:id/prs/:pr_number", repoCtrl.GetPR)

		// Actions
		api.POST("/repos/:id/prs/:pr_number/analyze", actionCtrl.AnalyzePR)
		api.POST("/repos/:id/release", actionCtrl.CreateRelease)

		// Metrics & Impact
		api.GET("/metrics", metricCtrl.GetMetrics)
		api.GET("/impact/:pr_id", actionCtrl.GetImpactAnalysis)
	}
}

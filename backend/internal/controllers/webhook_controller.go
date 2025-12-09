package controllers

import (
	"io"
	"net/http"

	"devplus-backend/internal/services"
	"devplus-backend/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type WebhookController struct {
	Service *services.WebhookService
}

func NewWebhookController() *WebhookController {
	return &WebhookController{
		Service: &services.WebhookService{},
	}
}

func (ctrl *WebhookController) HandleWebhook(c *gin.Context) {
	eventType := c.GetHeader("X-GitHub-Event")
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
		return
	}

	// Validate Signature
	if err := ctrl.Service.ValidateSignature(c.Request, body); err != nil {
		log.Warn().Err(err).Msg("Invalid webhook signature")
		utils.SendError(c, http.StatusUnauthorized, "Invalid signature", err)
		return
	}

	// Process in background to return 200 quickly to GitHub
	go func() {
		if err := ctrl.Service.ProcessWebhook(eventType, body); err != nil {
			log.Error().Err(err).Str("event", eventType).Msg("Failed to process webhook")
		} else {
			log.Info().Str("event", eventType).Msg("Webhook processed successfully")
		}
	}()

	c.JSON(http.StatusOK, gin.H{"status": "received"})
}

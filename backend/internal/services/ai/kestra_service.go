package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"

	"devplus-backend/internal/models"
)

type KestraAIService struct {
	kestraURL string
	username  string
	password  string
	client    *http.Client
}

func NewKestraAIService(kestraURL, username, password string) *KestraAIService {
	return &KestraAIService{
		kestraURL: kestraURL,
		username:  username,
		password:  password,
		client:    &http.Client{Timeout: 10 * time.Second},
	}
}

type KestraExecutionRequest struct {
	Namespace string                 `json:"namespace"`
	FlowId    string                 `json:"flowId"`
	Inputs    map[string]interface{} `json:"inputs"`
	Wait      bool                   `json:"wait"`
}

func (s *KestraAIService) AnalyzePR(ctx context.Context, pr *models.PullRequest, callbackURL string) error {
	log.Info().Str("pr_id", pr.ID).Str("flow_id", "ai-pull-request-analysis").Msg("[KestraService] Analyzing PR")

	// Ensure repository is loaded
	if pr.Repository == nil {
		log.Error().Str("pr_id", pr.ID).Msg("[KestraService] Repository not loaded in PR model")
		return fmt.Errorf("repository not loaded for PR %s", pr.ID)
	}

	// Construct inputs for Kestra Flow
	inputs := map[string]interface{}{
		"pr_id":        pr.ID,
		"pr_number":    pr.Number,
		"repo_id":      pr.RepoID,
		"repo_owner":   pr.Repository.Owner,
		"repo_name":    pr.Repository.Name,
		"pr_title":     pr.Title,
		"callback_url": callbackURL,
	}

	reqBody := KestraExecutionRequest{
		Namespace: "devplus",
		FlowId:    "ai-pull-request-analysis",
		Inputs:    inputs,
		Wait:      false,
	}

	// Use Webhook endpoint
	// Format: /api/v1/executions/webhook/{namespace}/{flowId}/{key}
	url := fmt.Sprintf("%s/api/v1/executions/webhook/%s/%s/devplus-webhook-key", s.kestraURL, reqBody.Namespace, reqBody.FlowId)

	// For Webhook, we send the inputs directly as the body, not wrapped in KestraExecutionRequest
	// So we marshal 'inputs' directly
	jsonBody, err := json.Marshal(inputs)

	if err != nil {
		log.Error().Err(err).Msg("[KestraService] Failed to marshal request body")
		return err
	}

	log.Debug().Str("url", url).RawJSON("body", jsonBody).Msg("[KestraService] Sending webhook request to Kestra")

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	if s.username != "" && s.password != "" {
		req.SetBasicAuth(s.username, s.password)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		log.Error().Err(err).Str("url", url).Msg("[KestraService] Failed to send request")
		return err
	}
	defer resp.Body.Close()

	log.Info().Int("status", resp.StatusCode).Msg("[KestraService] Received response from Kestra")

	if resp.StatusCode >= 300 {
		// Read body for error details
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		log.Error().Int("status", resp.StatusCode).Str("body", buf.String()).Msg("[KestraService] Kestra reported error")
		return fmt.Errorf("failed to trigger kestra workflow: %s | Body: %s", resp.Status, buf.String())
	}

	return nil
}

func (s *KestraAIService) AnalyzeRepo(ctx context.Context, repo *models.Repository, callbackURL string) error {
	log.Info().Str("repo_id", repo.ID).Str("flow_id", "ai-repo-analysis").Msg("[KestraService] Analyzing Repository")

	// Construct inputs for Kestra Flow
	inputs := map[string]interface{}{
		"repo_id":      repo.ID,
		"repo_name":    repo.Name,
		"callback_url": callbackURL,
	}

	reqBody := KestraExecutionRequest{
		Namespace: "devplus",
		FlowId:    "ai-repo-analysis",
		Inputs:    inputs,
		Wait:      false,
	}

	// Use Webhook endpoint
	// Format: /api/v1/executions/webhook/{namespace}/{flowId}/{key}
	url := fmt.Sprintf("%s/api/v1/executions/webhook/%s/%s/devplus-repo-webhook-key", s.kestraURL, reqBody.Namespace, reqBody.FlowId)

	// For Webhook, we send the inputs directly as the body, not wrapped in KestraExecutionRequest
	jsonBody, err := json.Marshal(inputs)
	if err != nil {
		log.Error().Err(err).Msg("[KestraService] Failed to marshal request body")
		return err
	}

	log.Debug().Str("url", url).RawJSON("body", jsonBody).Msg("[KestraService] Sending webhook request to Kestra")

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	if s.username != "" && s.password != "" {
		req.SetBasicAuth(s.username, s.password)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		log.Error().Err(err).Str("url", url).Msg("[KestraService] Failed to send request")
		return err
	}
	defer resp.Body.Close()

	log.Info().Int("status", resp.StatusCode).Msg("[KestraService] Received response from Kestra")

	if resp.StatusCode >= 300 {
		// Read body for error details
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		log.Error().Int("status", resp.StatusCode).Str("body", buf.String()).Msg("[KestraService] Kestra reported error")
		return fmt.Errorf("failed to trigger kestra workflow: %s | Body: %s", resp.Status, buf.String())
	}

	return nil
}

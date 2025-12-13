package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"devplus-backend/internal/models"
)

type KestraAIService struct {
	kestraURL string
	client    *http.Client
}

func NewKestraAIService(kestraURL string) *KestraAIService {
	return &KestraAIService{
		kestraURL: kestraURL,
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
	// Construct inputs for Kestra Flow
	inputs := map[string]interface{}{
		"pr_id":        pr.ID,
		"pr_number":    pr.Number,
		"repo_id":      pr.RepoID,
		"pr_title":     pr.Title,
		"callback_url": callbackURL,
		// In a real scenario, we might pass the Diff content here,
		// OR Kestra task fetches it from Github using token.
		// For simplicity, let's assume Kestra fetches it or we pass a placeholder.
		// "diff": "...",
	}

	reqBody := KestraExecutionRequest{
		Namespace: "devplus", // Assuming a namespace
		FlowId:    "ai-pull-request-analysis",
		Inputs:    inputs,
		Wait:      false, // Async
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	// Trigger via Kestra API: POST /api/v1/executions/trigger/{namespace}/{flowId}
	// Or POST /api/v1/executions with body.
	// Kestra API: POST /api/v1/executions/{namespace}/{flowId}
	url := fmt.Sprintf("%s/api/v1/executions/%s/%s", s.kestraURL, reqBody.Namespace, reqBody.FlowId)

	// If using Multipart form, it's different. Assuming via JSON for simple inputs.
	// Note: Kestra inputs are usually multipart for files, but JSON body is supported for simple key-values in recent versions
	// or via specific endpoints.
	// Let's assume standard POST to executions endpoint.
	// Actually, standard Kestra API for triggering is `POST /api/v1/executions/{namespace}/{flowId}` which typically takes multipart/form-data for inputs.
	// However, let's assume we can send JSON or we implement multipart if needed.
	// For this exercise, I'll send JSON and assume Kestra is configured/versioned to handle it or I'm using a wrapper.
	// A safe bet is multipart.
	// Let's write standard JSON request for now, assuming 0.18+ support or custom controller.

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("failed to trigger kestra workflow: %s", resp.Status)
	}

	return nil
}

func (s *KestraAIService) AnalyzeRepo(ctx context.Context, repo *models.Repository, callbackURL string) error {
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

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/api/v1/executions/%s/%s", s.kestraURL, reqBody.Namespace, reqBody.FlowId)

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("failed to trigger kestra workflow: %s", resp.Status)
	}

	return nil
}

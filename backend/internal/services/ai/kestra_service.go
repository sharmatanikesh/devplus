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

	// Fetch PR diff from GitHub API
	prDiff, err := s.fetchPRDiff(ctx, pr.Repository.Owner, pr.Repository.Name, int(*pr.Number))
	if err != nil {
		log.Error().Err(err).Str("pr_id", pr.ID).Msg("[KestraService] Failed to fetch PR diff")
		return fmt.Errorf("failed to fetch PR diff: %w", err)
	}

	// Construct inputs for Kestra Flow
	inputs := map[string]interface{}{
		"pr_id":        pr.ID,
		"pr_number":    pr.Number,
		"repo_id":      pr.RepoID,
		"repo_owner":   pr.Repository.Owner,
		"repo_name":    pr.Repository.Name,
		"pr_title":     pr.Title,
		"pr_diff":      prDiff,
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

	// Fetch README from GitHub API
	readme, err := s.fetchReadme(ctx, repo.Owner, repo.Name)
	if err != nil {
		log.Warn().Err(err).Str("repo_id", repo.ID).Msg("[KestraService] Failed to fetch README, using placeholder")
		readme = "No README available"
	}

	// Fetch file tree from GitHub API
	fileTree, err := s.fetchFileTree(ctx, repo.Owner, repo.Name)
	if err != nil {
		log.Warn().Err(err).Str("repo_id", repo.ID).Msg("[KestraService] Failed to fetch file tree, using placeholder")
		fileTree = "File tree not available"
	}

	// Construct inputs for Kestra Flow
	inputs := map[string]interface{}{
		"repo_id":      repo.ID,
		"repo_owner":   repo.Owner,
		"repo_name":    repo.Name,
		"readme":       readme,
		"file_tree":    fileTree,
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

// fetchReadme fetches the README content from GitHub API
func (s *KestraAIService) fetchReadme(ctx context.Context, owner, repo string) (string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/readme", owner, repo)
	
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Accept", "application/vnd.github.v3.raw")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	
	resp, err := s.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch readme: status %d", resp.StatusCode)
	}
	
	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(resp.Body); err != nil {
		return "", err
	}
	
	return buf.String(), nil
}

// fetchFileTree fetches the repository file tree from GitHub API
func (s *KestraAIService) fetchFileTree(ctx context.Context, owner, repo string) (string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/git/trees/main?recursive=1", owner, repo)
	
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	
	resp, err := s.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch file tree: status %d", resp.StatusCode)
	}
	
	var treeResponse struct {
		Tree []struct {
			Path string `json:"path"`
			Type string `json:"type"`
		} `json:"tree"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&treeResponse); err != nil {
		return "", err
	}
	
	// Build a simple text representation of the tree
	var fileTreeBuilder bytes.Buffer
	for _, item := range treeResponse.Tree {
		if item.Type == "blob" || item.Type == "tree" {
			fileTreeBuilder.WriteString(item.Path)
			fileTreeBuilder.WriteString("\n")
		}
	}
	
	return fileTreeBuilder.String(), nil
}

// fetchPRDiff fetches the pull request diff from GitHub API
func (s *KestraAIService) fetchPRDiff(ctx context.Context, owner, repo string, prNumber int) (string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/pulls/%d", owner, repo, prNumber)
	
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Accept", "application/vnd.github.v3.diff")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	
	resp, err := s.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch PR diff: status %d", resp.StatusCode)
	}
	
	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(resp.Body); err != nil {
		return "", err
	}
	
	return buf.String(), nil
}

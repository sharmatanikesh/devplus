package ai

import (
	"context"
	"devplus-backend/internal/models"
	"errors"
)

type AIService interface {
	AnalyzePR(ctx context.Context, pr *models.PullRequest, callbackURL string) error
	AnalyzeRepo(ctx context.Context, repo *models.Repository, callbackURL string) error
	TriggerReleaseRiskAnalysis(repoID string, owner string, name string, prData string, callbackURL string) error
}

type AIFactory struct {
	kestraURL      string
	kestraUsername string
	kestraPassword string
}

func NewAIFactory(kestraURL, kestraUsername, kestraPassword string) *AIFactory {
	return &AIFactory{
		kestraURL:      kestraURL,
		kestraUsername: kestraUsername,
		kestraPassword: kestraPassword,
	}
}

func (f *AIFactory) GetAIService(provider string) (AIService, error) {
	switch provider {
	case "kestra":
		return NewKestraAIService(f.kestraURL, f.kestraUsername, f.kestraPassword), nil
	default:
		return nil, errors.New("unsupported AI provider")
	}
}

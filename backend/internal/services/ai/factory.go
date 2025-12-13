package ai

import (
	"context"
	"devplus-backend/internal/models"
	"errors"
)

type AIService interface {
	AnalyzePR(ctx context.Context, pr *models.PullRequest, callbackURL string) error
	AnalyzeRepo(ctx context.Context, repo *models.Repository, callbackURL string) error
}

type AIFactory struct {
	kestraURL string
}

func NewAIFactory(kestraURL string) *AIFactory {
	return &AIFactory{
		kestraURL: kestraURL,
	}
}

func (f *AIFactory) GetAIService(provider string) (AIService, error) {
	switch provider {
	case "kestra":
		return NewKestraAIService(f.kestraURL), nil
	default:
		return nil, errors.New("unsupported AI provider")
	}
}

package service

import (
	"context"
	"exercise-generator/config"
	"exercise-generator/internal/model"
	"exercise-generator/internal/repository"

	_ "github.com/lib/pq"

	"go.uber.org/zap"
)

type Service struct {
	log        *zap.Logger
	cfg        *config.Config
	repository repository.IRepository
}

func NewService(logger *zap.Logger, cfg *config.Config) *Service {
	repo, err := repository.NewRepository(cfg)
	if err != nil {
		logger.Fatal("failed to init repository", zap.Error(err))
	}

	return &Service{
		log:        logger,
		cfg:        cfg,
		repository: repo,
	}
}

func (s *Service) CreateGenerationConfig(ctx context.Context, genCfg model.GenerationConfig) error {
	return s.repository.CreateGenerationConfig(ctx, repository.CreateGenerationConfigParams{
		Provider:         genCfg.Provider,
		GenModel:         genCfg.GenModel,
		TopP:             genCfg.TopP,
		Temperature:      genCfg.Temperature,
		MaxTokens:        genCfg.MaxTokens,
		SystemMessage:    genCfg.SystemMessage,
		AssistantMessage: genCfg.AssistantMessage,
		UserMessage:      genCfg.UserMessage,
	})
}

func (s *Service) UpdateGenerationConfig(ctx context.Context, genCfg model.GenerationConfig) error {
	return s.repository.UpdateGenerationConfig(ctx, repository.UpdateGenerationConfigParams{
		ID:               genCfg.ID,
		Provider:         genCfg.Provider,
		GenModel:         genCfg.GenModel,
		TopP:             genCfg.TopP,
		Temperature:      genCfg.Temperature,
		MaxTokens:        genCfg.MaxTokens,
		SystemMessage:    genCfg.SystemMessage,
		AssistantMessage: genCfg.AssistantMessage,
		UserMessage:      genCfg.UserMessage,
	})
}

func (s *Service) ListGenerationConfig(ctx context.Context) ([]model.GenerationConfig, error) {
	return s.repository.ListGenerationConfig(ctx)
}

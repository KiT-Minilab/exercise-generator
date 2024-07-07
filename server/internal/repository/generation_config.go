package repository

import (
	"context"
	"time"

	"exercise-generator/internal/model"
)

type IGenerationConfigRepository interface {
	CreateGenerationConfig(ctx context.Context, params CreateGenerationConfigParams) error
	UpdateGenerationConfig(ctx context.Context, params UpdateGenerationConfigParams) error
	ListGenerationConfig(ctx context.Context) ([]model.GenerationConfig, error)
	GetGenerationConfigByProvider(ctx context.Context, provider string) (*model.GenerationConfig, error)
}

type CreateGenerationConfigParams struct {
	Provider         string
	GenModel         string
	TopP             float32
	Temperature      float32
	MaxTokens        int
	SystemMessage    string
	AssistantMessage string
	UserMessage      string
}

func (r *repositoryImpl) CreateGenerationConfig(ctx context.Context, params CreateGenerationConfigParams) error {
	_, err := r.db.NamedExecContext(ctx, `
			INSERT INTO generation_config (provider, gen_model, top_p, temperature, max_tokens, system_message, assistant_message, user_message, created_at, updated_at)
			VALUES (:provider, :gen_model, :top_p, :temperature, :max_tokens, :system_message, :assistant_message, :user_message, :created_at, :updated_at)
		`,
		map[string]interface{}{
			"provider":          params.Provider,
			"gen_model":         params.GenModel,
			"top_p":             params.TopP,
			"temperature":       params.Temperature,
			"max_tokens":        params.MaxTokens,
			"system_message":    params.SystemMessage,
			"assistant_message": params.AssistantMessage,
			"user_message":      params.UserMessage,
			"created_at":        time.Now(),
			"updated_at":        time.Now(),
		})
	if err != nil {
		return err
	}
	return nil
}

type UpdateGenerationConfigParams struct {
	ID               int
	Provider         string
	GenModel         string
	TopP             float32
	Temperature      float32
	MaxTokens        int
	SystemMessage    string
	AssistantMessage string
	UserMessage      string
}

func (r *repositoryImpl) UpdateGenerationConfig(ctx context.Context, params UpdateGenerationConfigParams) error {
	_, err := r.db.NamedExecContext(ctx, `
			UPDATE generation_config
		 	SET provider = :provider, gen_model = :gen_model, top_p = :top_p, temperature = :temperature, max_tokens = :max_tokens, system_message = :system_message, assistant_message = :assistant_message, user_message = :user_message, updated_at = :updated_at
			WHERE id = :id`,
		map[string]interface{}{
			"id":                params.ID,
			"provider":          params.Provider,
			"gen_model":         params.GenModel,
			"top_p":             params.TopP,
			"temperature":       params.Temperature,
			"max_tokens":        params.MaxTokens,
			"system_message":    params.SystemMessage,
			"assistant_message": params.AssistantMessage,
			"user_message":      params.UserMessage,
			"updated_at":        time.Now(),
		})
	if err != nil {
		return err
	}
	return nil
}

func (r *repositoryImpl) ListGenerationConfig(ctx context.Context) ([]model.GenerationConfig, error) {
	var genCfgs []model.GenerationConfig
	err := r.db.SelectContext(ctx, &genCfgs, "SELECT * FROM generation_config")
	if err != nil {
		return nil, err
	}
	return genCfgs, nil
}

func (r *repositoryImpl) GetGenerationConfigByProvider(ctx context.Context, provider string) (*model.GenerationConfig, error) {
	var genCfg model.GenerationConfig
	err := r.db.GetContext(ctx, &genCfg, "SELECT * FROM generation_config WHERE provider = $1", provider)
	if err != nil {
		return nil, err
	}
	return &genCfg, nil
}

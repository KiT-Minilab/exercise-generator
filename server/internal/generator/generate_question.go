package generator

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"exercise-generator/config"
	// "exercise-generator/internal/adapter/gemini"
	"exercise-generator/internal/adapter/openai"
	"exercise-generator/internal/constant"
	"exercise-generator/internal/model"
	"exercise-generator/internal/repository"

	"go.uber.org/zap"
)

type QuestionGenerator interface {
	GenerateEnglishMultipleChoicesQuestion(ctx context.Context, word string, questionType string) (*model.MultipleChoicesQuestion, error)
}

type questionGenerator struct {
	cfg        *config.Config
	logger     *zap.Logger
	provider   string
	repository repository.IRepository
}

func NewQuestionGenerator(cfg *config.Config, logger *zap.Logger, provider string, repository repository.IRepository) QuestionGenerator {
	return &questionGenerator{
		cfg:        cfg,
		logger:     logger,
		provider:   provider,
		repository: repository,
	}
}

func (q *questionGenerator) GenerateEnglishMultipleChoicesQuestion(ctx context.Context, word string, questionType string) (*model.MultipleChoicesQuestion, error) {
	switch strings.ToUpper(q.provider) {
	case constant.ProviderOpenAi:
		openaiAdapter := openai.NewOpenaiAdapter(q.cfg, q.logger)
		genConfig, err := q.repository.GetGenerationConfigByProvider(ctx, constant.ProviderOpenAi)
		if err != nil {
			return nil, err
		}
		return openaiAdapter.GenerateEnglishMultipleChoicesQuestion(ctx, *genConfig, model.QuestionRequest{
			Word:         word,
			QuestionType: questionType,
		})
	case constant.ProviderGoogle:
		return nil, errors.New("current not work with Gemini model")
		// return gemini.NewGeminiAdapter(cfg, logger), nil
	default:
		return nil, fmt.Errorf("unsupported provider: %s", q.provider)
	}

}

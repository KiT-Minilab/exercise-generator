package gemini

import (
	"context"
	"exercise-generator/config"
	"exercise-generator/internal/model"

	"go.uber.org/zap"
)

type GeminiAdapter interface {
	GenerateEnglishMultipleChoicesQuestion(ctx context.Context, word string, questionType string) (*model.MultipleChoicesQuestion, error)
}

type geminiAdapterImpl struct {
}

func NewGeminiAdapter(cfg *config.Config, logger *zap.Logger) GeminiAdapter {
	return &geminiAdapterImpl{}
}

func (g *geminiAdapterImpl) GenerateEnglishMultipleChoicesQuestion(ctx context.Context, word string, questionType string) (*model.MultipleChoicesQuestion, error) {
	return &model.MultipleChoicesQuestion{}, nil
}

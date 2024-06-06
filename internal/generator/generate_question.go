package generator

import (
	"context"
	"errors"
	"fmt"

	"exercise-generator/config"
	// "exercise-generator/internal/adapter/gemini"
	"exercise-generator/internal/adapter/openai"
	"exercise-generator/internal/constant"
	"exercise-generator/internal/model"

	"go.uber.org/zap"
)

type QuestionGenerator interface {
	GenerateEnglishMultipleChoicesQuestion(ctx context.Context, word string, questionType string) (*model.MultipleChoicesQuestion, error)
}

func NewQuestionGenerator(cfg *config.Config, logger *zap.Logger, model model.AIModel) (QuestionGenerator, error) {
	switch model {
	case constant.ModelOpenAI:
		return openai.NewOpenaiAdapter(cfg, logger), nil
	case constant.ModelGemini:
		return nil, errors.New("current not work with Gemini model")
		// return gemini.NewGeminiAdapter(cfg, logger), nil
	default:
		return nil, fmt.Errorf("unsupported model: %s", model)
	}
}

package openai

import (
	"context"
	"encoding/json"
	"fmt"

	"exercise-generator/config"
	"exercise-generator/internal/model"

	"github.com/sashabaranov/go-openai"
	"go.uber.org/zap"
)

type OpenaiAdapter interface {
	GenerateEnglishMultipleChoicesQuestion(ctx context.Context, word string, questionType string) (*model.MultipleChoicesQuestion, error)
}

type openaiAdapterImpl struct {
	cfg    *config.Config
	client *openai.Client
	log    *zap.Logger
}

func NewOpenaiAdapter(cfg *config.Config, logger *zap.Logger) OpenaiAdapter {
	openaiClient := openai.NewClient(cfg.OpenaiClient.ApiKey)

	return &openaiAdapterImpl{
		cfg:    cfg,
		client: openaiClient,
		log:    logger,
	}
}

func (o *openaiAdapterImpl) GenerateEnglishMultipleChoicesQuestion(ctx context.Context, word string, questionType string) (*model.MultipleChoicesQuestion, error) {
	logger := o.log.Named("openaiAdapterImpl.GenerateQuestion").With(zap.String("word", word), zap.String("question_type", questionType))

	resp, err := o.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: openai.GPT3Dot5Turbo0125,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: o.cfg.OpenaiClient.SystemMessage,
			},
			{
				Role:    openai.ChatMessageRoleAssistant,
				Content: o.cfg.OpenaiClient.AssistanceMessage,
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: fmt.Sprintf(o.cfg.OpenaiClient.UserMessageTemplate, questionType, word, word),
			},
		},
		MaxTokens:   o.cfg.OpenaiClient.MaxTokens,
		Temperature: o.cfg.OpenaiClient.Temperature,
		TopP:        o.cfg.OpenaiClient.TopP,
		ResponseFormat: &openai.ChatCompletionResponseFormat{
			Type: openai.ChatCompletionResponseFormatTypeJSONObject,
		},
	})
	if err != nil {
		logger.Error("failed to create chat completion", zap.Error(err))
		return nil, err
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("choices are empty")
	}
	logger.Info("successfully generated question", zap.String("response", resp.Choices[0].Message.Content))

	var question model.MultipleChoicesQuestion
	err = json.Unmarshal([]byte(resp.Choices[0].Message.Content), &question)
	if err != nil {
		logger.Error("failed to unmarshal response", zap.Error(err))
	}

	return &question, nil
}

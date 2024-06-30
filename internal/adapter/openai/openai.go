package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"text/template"

	"exercise-generator/config"
	"exercise-generator/internal/model"

	"github.com/sashabaranov/go-openai"
	"go.uber.org/zap"
)

type OpenaiAdapter interface {
	GenerateEnglishMultipleChoicesQuestion(ctx context.Context, genConfig model.GenerationConfig, questionRequest model.QuestionRequest) (*model.MultipleChoicesQuestion, error)
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

func (o *openaiAdapterImpl) GenerateEnglishMultipleChoicesQuestion(ctx context.Context, genCfg model.GenerationConfig, questionRequest model.QuestionRequest) (*model.MultipleChoicesQuestion, error) {
	logger := o.log.Named("openaiAdapterImpl.GenerateQuestion")

	t := template.Must(template.New("userMessage").Parse(genCfg.UserMessage))

	var userMessageBuf bytes.Buffer

	// Execute the template and write to the buffer
	err := t.Execute(&userMessageBuf, questionRequest)
	if err != nil {
		logger.Error("failed to execute template", zap.Error(err))
		return nil, err
	}

	resp, err := o.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: genCfg.GenModel,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: genCfg.SystemMessage,
			},
			{
				Role:    openai.ChatMessageRoleAssistant,
				Content: genCfg.AssistantMessage,
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: userMessageBuf.String(),
			},
		},
		MaxTokens:   genCfg.MaxTokens,
		Temperature: genCfg.Temperature,
		TopP:        genCfg.TopP,
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

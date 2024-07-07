package gemini

// import (
// 	"context"
// 	"encoding/json"
// 	"exercise-generator/config"
// 	"exercise-generator/internal/model"
// 	"fmt"

// 	"github.com/google/generative-ai-go/genai"
// 	"go.uber.org/zap"
// 	"google.golang.org/api/option"
// )

// type GeminiAdapter interface {
// 	GenerateEnglishMultipleChoicesQuestion(ctx context.Context, word string, questionType string) (*model.MultipleChoicesQuestion, error)
// }

// type geminiAdapterImpl struct {
// 	cfg    *config.Config
// 	client *genai.Client
// 	log    *zap.Logger
// }

// func NewGeminiAdapter(cfg *config.Config, logger *zap.Logger) GeminiAdapter {
// 	ctx := context.Background()
// 	client, err := genai.NewClient(ctx, option.WithAPIKey(cfg.GeminiClient.ApiKey))
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer client.Close()

// 	return &geminiAdapterImpl{client: client, cfg: cfg, log: logger}
// }

// func (g *geminiAdapterImpl) GenerateEnglishMultipleChoicesQuestion(ctx context.Context, word string, questionType string) (*model.MultipleChoicesQuestion, error) {
// 	logger := g.log.Named("geminiAdapterImpl.GenerateQuestion").With(zap.String("word", word), zap.String("question_type", questionType))

// 	genModel := g.client.GenerativeModel(g.cfg.GeminiClient.Model)
// 	genModel.SetTemperature(g.cfg.GeminiClient.Temperature)
// 	genModel.SetTopP(g.cfg.GeminiClient.TopP)
// 	genModel.SetMaxOutputTokens(int32(g.cfg.GeminiClient.MaxTokens))

// 	promptParts := make([]genai.Part, 0)
// 	for _, part := range g.cfg.GeminiClient.PromptParts {
// 		promptParts = append(promptParts, genai.Text(part))
// 	}

// 	// insert your request at last
// 	promptParts = append(promptParts,
// 		[]genai.Part{
// 			genai.Text(fmt.Sprintf("word: %s", word)),
// 			genai.Text(fmt.Sprintf("questionType: %s", questionType)),
// 			genai.Text("output: "),
// 		}...)

// 	resp, err := genModel.GenerateContent(ctx, promptParts...)
// 	if err != nil {
// 		logger.Error("Failed to generate content", zap.Error(err))
// 		return nil, err
// 	}
// 	logger.Info("Generated content", zap.Reflect("resp", resp))

// 	content := ""
// 	for _, part := range resp.Candidates[0].Content.Parts {
// 		textPart, ok := part.(genai.Text)
// 		if !ok {
// 			return nil, fmt.Errorf("part is not of type genai.Text")
// 		}
// 		content += fmt.Sprintf("%s\n", string(textPart))
// 	}
// 	logger.Info("Content", zap.String("content", content))

// 	var multipleChoicesQuestion *model.MultipleChoicesQuestion
// 	err = json.Unmarshal([]byte(content), &multipleChoicesQuestion)
// 	if err != nil {
// 		logger.Error("Failed to unmarshal content", zap.Error(err))
// 		return nil, err
// 	}
// 	return multipleChoicesQuestion, nil
// }

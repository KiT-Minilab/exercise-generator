package main

import (
	"context"
	"encoding/json"
	"errors"
	"exercise-generator/config"
	"exercise-generator/script"
	"flag"
	"fmt"
	"log"
	"os"

	openai "github.com/sashabaranov/go-openai"
	"go.uber.org/zap"
)

const (
	QuestionTypeDefinition  = "definition"
	QuestionTypeSynonym     = "synonym"
	QuestionTypeApplication = "application"
)

func main() {

	// Define command line flags
	var (
		scriptName   string
		questionType string
		argWord      string
	)
	flag.StringVar(&scriptName, "script", "", "Enter your script name")
	flag.StringVar(&questionType, "questionType", "", "Enter your question type")
	flag.StringVar(&argWord, "word", "", "Enter your word")

	// Parse command line flags
	flag.Parse()

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	cfg, err := config.Load()
	if err != nil {
		logger.Fatal("failed to load config", zap.Error(err))
	}

	openaiClient := openai.NewClient(cfg.OpenaiClient.ApiKey)
	switch scriptName {
	case "baseline":
		err = handleScriptBaseline()
	case "generate-question":
		if questionType != QuestionTypeDefinition && questionType != QuestionTypeSynonym && questionType != QuestionTypeApplication {
			logger.Fatal("Invalid question type", zap.String("questionType", questionType))
		}

		question, err := generateQuestion(openaiClient, logger, argWord, questionType)
		if err != nil {
			logger.Error("failed to generate question", zap.Error(err), zap.String("word", argWord))
		}
		logger.Info("successfully generated question", zap.Reflect("question", question))
	case "generate-baseline-question":
		words, err := script.ReadBaselineWords()
		if err != nil {
			logger.Fatal("failed to read baseline words", zap.Error(err))
		}

		if questionType != QuestionTypeDefinition && questionType != QuestionTypeSynonym && questionType != QuestionTypeApplication {
			logger.Fatal("Invalid question type", zap.String("questionType", questionType))
		}

		questions := make([]Question, 0, len(words))
		for index, word := range words {
			question, err := generateQuestion(openaiClient, logger, word, questionType)
			if err != nil {
				logger.Error("failed to generate question", zap.Error(err), zap.String("word", word))
			}
			logger.Info("successfully generated question", zap.Reflect("question", question), zap.Int("index", index))
			questions = append(questions, question)
		}

		// Marshal the array of JSON objects into JSON format
		jsonData, err := json.MarshalIndent(questions, "", "    ")
		if err != nil {
			logger.Error("Error marshaling JSON:", zap.Error(err))
			return
		}

		file, err := os.Create(fmt.Sprintf("data/v2/%s.json", questionType))
		defer file.Close()
		if err != nil {
			logger.Error("failed to create definition.json", zap.Error(err))
			return
		}

		_, err = file.Write(jsonData)
		if err != nil {
			logger.Error("failed to write json file", zap.Error(err))
			return
		}
	default:
		return
	}

	if err != nil {
		log.Fatalln("Error:", err)
	}
}

func handleScriptBaseline() error {
	words, err := script.ReadBaselineWords()
	if err != nil {
		return err
	}
	fmt.Println(words)
	return nil
}

type Question struct {
	QuestionStem string   `json:"questionStem"`
	QuestionType string   `json:"questionType"`
	Answers      []string `json:"answers"`
	Distractors  []string `json:"distractors"`
}

func generateQuestion(client *openai.Client, logger *zap.Logger, word string, questionType string) (Question, error) {
	systemMessage := "You are an English tutor who teach English for ESL learner"
	assistanceMessage := "Need to give student questions to learn vocabulary. There are 3 types: Definition, Synonym and Application. Here are some examples for each type in JSON format, the example is wrapped by triple backtick (```):\n1. \"Definition\" question: ```{\"questionStem\":\"**cone** means\",    \"questionType\": \"definition\",\"answers\": [\"a shape with a circular base and sides tapering to a point\"],\"distractors\":[\"audacious behavior that you have no right to\",\"a wide scope\",\"a quality that arouses emotions\"]}```\n2. \"Synonym\": question:```{\"questionStem\":\"**diminish** has the same or almost the same meaning as\",\"questionType\":\"synonym\",\"answers\":[\"fall\"],\"distractors\":[\"restore\",\"quiver\",\"glare\"]}```\n3. \"Application\" question:```{\"questionStem\":\"Which of the following would most likely **glint**?\",\"questionType\":\"application\",\"answers\":[\"a newly polished car in the sunlight\"],\"distractors\":[\"a newly written manuscript on an editor’s desk\",\"an old sweater in a trunk in the attic\",\"a fine powder used to absorb moisture\"]}```. Note that \"Application\" questions might not suitable for all words, especially for abstract words (e.g: \"apparently\", \"respective\",...)."
	userMessage := fmt.Sprintf("Give a single choice question with %s type to help student learn word **%s**. Provide the question in JSON format with keys: `questionStem`, `questionType`, `answers`, `distractors`. Note that, there must be only one correct answer (length of `answers` array is 1), `distractors` must be 3 wrong choices, the word **%s** must only appear in field `questionStem` and not appear in 2 fields `answers` and `distractors`. If there is no sensible question for the word, let `answers` and `distractors` empty", questionType, word, word)

	var resp openai.ChatCompletionResponse
	err := errors.New("error")
	retryCount := 0
	for err != nil && retryCount <= 5 {
		resp, err = client.CreateChatCompletion(context.Background(), openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo0125,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: systemMessage,
				},
				{
					Role:    openai.ChatMessageRoleAssistant,
					Content: assistanceMessage,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: userMessage,
				},
			},
			MaxTokens:   256,
			Temperature: 1,
			TopP:        1,
			ResponseFormat: &openai.ChatCompletionResponseFormat{
				Type: openai.ChatCompletionResponseFormatTypeJSONObject,
			},
		})
	}

	if len(resp.Choices) == 0 {
		err := errors.New("choices are empty")
		logger.Error("there is no choice in response", zap.Error(err))
		return Question{}, err
	}
	logger.Info("create chat completion successfully ", zap.Reflect("choices", resp.Choices))

	var question Question
	err = json.Unmarshal([]byte(resp.Choices[0].Message.Content), &question)
	if err != nil {
		logger.Error("failed to unmarshal response", zap.Error(err))
	}
	return question, nil
}

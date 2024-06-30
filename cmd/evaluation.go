package main

import (
	"context"
	"exercise-generator/internal/constant"
	"exercise-generator/model_evaluator"
	"fmt"
	"strings"

	"github.com/urfave/cli/v3"
)

func validateQuestionType(cmdQuestionType string) ([]string, error) {
	if cmdQuestionType == "" {
		return []string{constant.QuestionTypeDefinition, constant.QuestionTypeSynonym, constant.QuestionTypeApplication}, nil
	}

	if strings.EqualFold(cmdQuestionType, constant.QuestionTypeDefinition) ||
		strings.EqualFold(cmdQuestionType, constant.QuestionTypeSynonym) ||
		strings.EqualFold(cmdQuestionType, constant.QuestionTypeApplication) {
		return []string{cmdQuestionType}, nil
	}

	return nil, fmt.Errorf("invalid question type: %s", cmdQuestionType)
}

func evaluateBaselineWords(ctx context.Context, cmd *cli.Command) error {
	var (
		cmdQuestionType = cmd.String("questionType")
		cmdProvider     = cmd.String("provider")
	)

	questionTypes, err := validateQuestionType(cmdQuestionType)
	if err != nil {
		return err
	}

	if !strings.EqualFold(cmdProvider, string(constant.ProviderOpenAi)) &&
		!strings.EqualFold(cmdProvider, string(constant.ProviderGoogle)) {
		err := fmt.Errorf("invalid provider: %s", cmdProvider)
		return err
	}

	evaluator := model_evaluator.NewModelEvaluator(cfg, logger)
	return evaluator.EvaluateBaselineWords(questionTypes, cmdProvider)
}

// func generateBaselineQuestions(ctx context.Context, cmd *cli.Command) error {
// 	var (
// 		cmdQuestionType = cmd.String("questionType")
// 		cmdProvider     = cmd.String("provider")
// 	)

// 	questionTypes, err := validateQuestionType(cmdQuestionType)
// 	if err != nil {
// 		return err
// 	}

// 	if cmdProvider != string(constant.ProviderOpenAi) && cmdProvider != string(constant.ProviderGoogle) {
// 		err := fmt.Errorf("invalid provider: %s", cmdProvider)
// 		return err
// 	}

// 	content, err := os.ReadFile("script/baseline_words.txt")
// 	if err != nil {
// 		return err
// 	}

// 	// Split the content into words
// 	words := strings.Fields(string(content))
// 	if err != nil {
// 		return err
// 	}

// 	repo, err := repository.NewRepository(cfg)
// 	if err != nil {
// 		return err
// 	}
// 	questionGenerator := generator.NewQuestionGenerator(cfg, logger, cmdProvider, repo)

// 	file, err := os.Create(fmt.Sprintf("%s/eval_%d.csv", cfg.EvaluationFilePath, time.Now().Unix()))
// 	if err != nil {
// 		return err
// 	}
// 	defer file.Close()

// 	writer := csv.NewWriter(file)
// 	defer writer.Flush()

// 	// write the header
// 	writer.Write([]string{"question", "question type", "answer", "distractor 1", "distractor 2", "distractor 3"})

// 	for _, questionType := range questionTypes {
// 		for index, word := range words {
// 			question, err := questionGenerator.GenerateEnglishMultipleChoicesQuestion(context.Background(), word, questionType)
// 			if err != nil {
// 				logger.Error("failed to generate question", zap.Error(err), zap.String("word", word))
// 			}
// 			logger.Info("successfully generated question", zap.Reflect("question", question), zap.Int("index", index))

// 			err = writeQuestionToFile(writer, question)
// 			if err != nil {
// 				logger.Error("failed to write question to file", zap.Error(err))
// 			}
// 			time.Sleep(time.Duration(cfg.EvaluationIntervalSeconds) * time.Second)
// 		}
// 	}

// 	return nil
// }

// func writeQuestionToFile(writer *csv.Writer, question *model.MultipleChoicesQuestion) error {
// 	if question == nil {
// 		return nil
// 	}

// 	var (
// 		answer      string
// 		distractor1 string
// 		distractor2 string
// 		distractor3 string
// 	)

// 	if question.Answers != nil && len(question.Answers) >= 0 {
// 		answer = question.Answers[0]
// 	}

// 	if question.Distractors != nil && len(question.Distractors) >= 3 {
// 		distractor1 = question.Distractors[0]
// 		distractor2 = question.Distractors[1]
// 		distractor3 = question.Distractors[2]
// 	}

// 	return writer.Write([]string{question.QuestionStem, question.QuestionType, answer, distractor1, distractor2, distractor3})
// }

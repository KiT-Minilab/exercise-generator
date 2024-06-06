package main

import (
	"context"
	"encoding/csv"
	"exercise-generator/internal/constant"
	"exercise-generator/internal/generator"
	"exercise-generator/internal/model"
	"exercise-generator/script"
	"fmt"
	"os"
	"time"

	"github.com/urfave/cli/v3"
	"go.uber.org/zap"
)

func generateBaselineQuestions(ctx context.Context, cmd *cli.Command) error {
	var (
		cmdQuestionType = cmd.String("questionType")
		cmdAiModel      = cmd.String("aiModel")
	)

	questionTypes, err := validateQuestionType(cmdQuestionType)
	if err != nil {
		return err
	}

	if cmdAiModel != string(constant.ModelOpenAI) && cmdAiModel != string(constant.ModelGemini) {
		err := fmt.Errorf("invalid AI model: %s", cmdAiModel)
		return err
	}

	words, err := script.ReadBaselineWords()
	if err != nil {
		return err
	}

	questionGenerator, err := generator.NewQuestionGenerator(cfg, logger, model.AIModel(cmdAiModel))
	if err != nil {
		return err
	}

	file, err := os.Create(fmt.Sprintf("%s/eval_%d.csv", cfg.EvaluationFilePath, time.Now().Unix()))
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// write the header
	writer.Write([]string{"question", "question type", "answer", "distractor 1", "distractor 2", "distractor 3"})

	for _, questionType := range questionTypes {
		for index, word := range words {
			question, err := questionGenerator.GenerateEnglishMultipleChoicesQuestion(context.Background(), word, questionType)
			if err != nil {
				logger.Error("failed to generate question", zap.Error(err), zap.String("word", word))
			}
			logger.Info("successfully generated question", zap.Reflect("question", question), zap.Int("index", index))

			err = writeQuestionToFile(writer, question)
			if err != nil {
				logger.Error("failed to write question to file", zap.Error(err))
			}
			time.Sleep(time.Duration(cfg.EvaluationIntervalSeconds) * time.Second)
		}
	}

	return nil
}

func validateQuestionType(cmdQuestionType string) ([]string, error) {
	if cmdQuestionType == "" {
		return []string{constant.QuestionTypeDefinition, constant.QuestionTypeSynonym, constant.QuestionTypeApplication}, nil
	}

	if cmdQuestionType == constant.QuestionTypeDefinition || cmdQuestionType == constant.QuestionTypeSynonym || cmdQuestionType == constant.QuestionTypeApplication {
		return []string{cmdQuestionType}, nil
	}

	return nil, fmt.Errorf("invalid question type: %s", cmdQuestionType)
}

func writeQuestionToFile(writer *csv.Writer, question *model.MultipleChoicesQuestion) error {
	if question == nil {
		return nil
	}

	var (
		answer      string
		distractor1 string
		distractor2 string
		distractor3 string
	)

	if question.Answers != nil && len(question.Answers) >= 0 {
		answer = question.Answers[0]
	}

	if question.Distractors != nil && len(question.Distractors) >= 3 {
		distractor1 = question.Distractors[0]
		distractor2 = question.Distractors[1]
		distractor3 = question.Distractors[2]
	}

	return writer.Write([]string{question.QuestionStem, question.QuestionType, answer, distractor1, distractor2, distractor3})
}

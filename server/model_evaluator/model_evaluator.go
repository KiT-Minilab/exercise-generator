package model_evaluator

import (
	"context"
	"encoding/csv"
	"exercise-generator/config"
	"exercise-generator/internal/generator"
	"exercise-generator/internal/model"
	"exercise-generator/internal/repository"
	"fmt"
	"os"
	"strings"
	"time"

	"go.uber.org/zap"
)

type ModelEvaluator struct {
	cfg    *config.Config
	logger *zap.Logger
}

func NewModelEvaluator(cfg *config.Config, logger *zap.Logger) *ModelEvaluator {
	return &ModelEvaluator{cfg: cfg, logger: logger}
}

func (e *ModelEvaluator) EvaluateBaselineWords(questionTypes []string, provider string) error {
	baseLineWordPath := "./model_evaluator/baseline_words.txt"
	if e.cfg.EvaluationConfig.BaselineFilePath != "" {
		baseLineWordPath = e.cfg.EvaluationConfig.BaselineFilePath
	}

	content, err := os.ReadFile(baseLineWordPath)
	if err != nil {
		return err
	}

	// Split the content into words
	words := strings.Fields(string(content))

	repo, err := repository.NewRepository(e.cfg)
	if err != nil {
		return err
	}
	questionGenerator := generator.NewQuestionGenerator(e.cfg, e.logger, provider, repo)

	resultDir := "./result"
	if e.cfg.EvaluationConfig.ResultDir != "" {
		resultDir = e.cfg.EvaluationConfig.ResultDir
	}
	file, err := os.Create(fmt.Sprintf("%s/eval_%d.csv", resultDir, time.Now().Unix()))
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
				e.logger.Error("failed to generate question", zap.Error(err), zap.String("word", word))
				continue
			}
			e.logger.Info("successfully generated question", zap.Reflect("question", question), zap.Int("index", index))

			err = e.writeQuestionToFile(writer, question)
			if err != nil {
				e.logger.Error("failed to write question to file", zap.Error(err))
			}
			time.Sleep(time.Duration(e.cfg.EvaluationConfig.IntervalSeconds) * time.Second)
		}
	}

	return nil
}

func (e *ModelEvaluator) writeQuestionToFile(writer *csv.Writer, question *model.MultipleChoicesQuestion) error {
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

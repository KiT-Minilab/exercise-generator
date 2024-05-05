package model

type MultipleChoicesQuestion struct {
	QuestionStem string   `json:"questionStem"`
	QuestionType string   `json:"questionType"`
	Answers      []string `json:"answers"`
	Distractors  []string `json:"distractors"`
}

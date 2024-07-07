package model

import (
	"time"
)

type GenerationConfig struct {
	ID               int       `json:"id" db:"id"`
	Provider         string    `json:"provider" db:"provider"`
	GenModel         string    `json:"genModel" db:"gen_model"`
	TopP             float32   `json:"topP" db:"top_p"`
	Temperature      float32   `json:"temperature" db:"temperature"`
	MaxTokens        int       `json:"maxTokens" db:"max_tokens"`
	SystemMessage    string    `json:"systemMessage" db:"system_message"`
	AssistantMessage string    `json:"assistantMessage" db:"assistant_message"`
	UserMessage      string    `json:"userMessage" db:"user_message"`
	CreatedAt        time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt        time.Time `json:"updatedAt" db:"updated_at"`
}

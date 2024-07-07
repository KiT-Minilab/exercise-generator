package config

// Config application
type Config struct {
	PostgreSQL DBConfig `json:"postgresql" mapstructure:"postgresql"`

	OpenaiClient OpenaiClient `json:"openai_client" mapstructure:"openai_client"`
	GeminiClient GeminiClient `json:"gemini_client" mapstructure:"gemini_client"`

	MigrationFolder string `json:"migration_folder" mapstructure:"migration_folder"`

	EvaluationConfig EvaluationConfig `json:"evaluation_config" mapstructure:"evaluation_config"`
}

type OpenaiClient struct {
	Host                string  `json:"host" mapstructure:"host"`
	ApiKey              string  `json:"api_key" mapstructure:"api_key"`
	Model               string  `json:"model" mapstructure:"model"`
	SystemMessage       string  `json:"system_message" mapstructure:"system_message"`
	AssistanceMessage   string  `json:"assistance_message" mapstructure:"assistance_message"`
	UserMessageTemplate string  `json:"user_message_template" mapstructure:"user_message_template"`
	MaxTokens           int     `json:"max_tokens" mapstructure:"max_tokens"`
	Temperature         float32 `json:"temperature" mapstructure:"temperature"`
	TopP                float32 `json:"top_p" mapstructure:"top_p"`
}

type GeminiClient struct {
	Host        string   `json:"host" mapstructure:"host"`
	ApiKey      string   `json:"api_key" mapstructure:"api_key"`
	Model       string   `json:"model" mapstructure:"model"`
	MaxTokens   int      `json:"max_tokens" mapstructure:"max_tokens"`
	PromptParts []string `json:"prompt_parts" mapstructure:"prompt_parts"`
	Temperature float32  `json:"temperature" mapstructure:"temperature"`
	TopP        float32  `json:"top_p" mapstructure:"top_p"`
}

type EvaluationConfig struct {
	ResultDir        string `json:"result_dir" mapstructure:"result_dir"`
	BaselineFilePath string `json:"baseline_file_path" mapstructure:"baseline_file_path"`
	IntervalSeconds  int    `json:"interval_seconds" mapstructure:"interval_seconds"`
}

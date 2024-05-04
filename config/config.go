package config

// Config application
type Config struct {
	OpenaiClient OpenaiClient `json:"openai_client"`
}

type OpenaiClient struct {
	Host   string `json:"host"`
	ApiKey string `json:"api_key"`
}

func Load() (*Config, error) {
	// TODO: load config from env variables
	return loadDefaultConfig(), nil
}

// nolint
func loadDefaultConfig() *Config {
	c := &Config{
		OpenaiClient: OpenaiClient{
			Host:   "https://api.openai.com",
			ApiKey: "sk-proj-MpqloVzjiROlyTsnqJfTT3BlbkFJN52AbvMZHwefwQqQxMXI",
		},
	}
	return c
}

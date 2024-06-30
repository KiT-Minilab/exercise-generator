package config

import (
	"bytes"
	"encoding/json"
	"log"
	"strings"

	"github.com/spf13/viper"
)

func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/KiT/exercise-generator")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "__"))
	viper.AutomaticEnv()

	/**
	|-------------------------------------------------------------------------
	| You should set default config value here
	| 1. Populate the default value in (Source code)
	| 2. Then merge from config (YAML) and OS environment
	|-----------------------------------------------------------------------*/
	c := loadDefaultConfig()
	if configBuffer, err := json.Marshal(c); err != nil {
		log.Println("Oops! Marshal config is failed. ", err)
		return nil, err
	} else if err := viper.ReadConfig(bytes.NewBuffer(configBuffer)); err != nil {
		log.Println("Oops! Read default config is failed. ", err)
		return nil, err
	}
	if err := viper.MergeInConfig(); err != nil {
		log.Println("Read config file failed.", err)
	}

	// Populate the config again
	err := viper.Unmarshal(c)
	return c, err
}

// nolint
func loadDefaultConfig() *Config {
	return &Config{
		PostgreSQL:      PostgresSQLDefaultConfig(),
		MigrationFolder: "file://migrations",
		OpenaiClient: OpenaiClient{
			Host:   "https://api.openai.com",
			ApiKey: "",
		},
		GeminiClient: GeminiClient{
			Model: "gemini-1.0-pro",
		},
		EvaluationConfig: EvaluationConfig{
			IntervalSeconds: 1,
		},
	}
}

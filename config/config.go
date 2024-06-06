package config

import (
	"bytes"
	"encoding/json"
	"log"
	"strings"

	"github.com/spf13/viper"
)

// Config application
type Config struct {
	PostgreSQL DBConfig `json:"postgresql"`

	OpenaiClient OpenaiClient `json:"openai_client"`
	GeminiClient GeminiClient `json:"gemini_client"`

	MigrationFolder           string `json:"migration_folder"`
	DataFolder                string `json:"data_folder"`
	EvaluationFilePath        string `json:"evaluation_file_path"`
	EvaluationIntervalSeconds int    `json:"evaluation_interval_seconds"`
}

type OpenaiClient struct {
	Host                string  `json:"host"`
	ApiKey              string  `json:"api_key"`
	Model               string  `json:"model"`
	SystemMessage       string  `json:"system_message"`
	AssistanceMessage   string  `json:"assistance_message"`
	UserMessageTemplate string  `json:"user_message_template"`
	MaxTokens           int     `json:"max_tokens"`
	Temperature         float32 `json:"temperature"`
	TopP                float32 `json:"top_p"`
}

type GeminiClient struct {
	Host        string   `json:"host"`
	ApiKey      string   `json:"api_key"`
	Model       string   `json:"model"`
	MaxTokens   int      `json:"max_tokens"`
	PromptParts []string `json:"prompt_parts"`
	Temperature float32  `json:"temperature"`
	TopP        float32  `json:"top_p"`
}

func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
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
			Host:                "https://api.openai.com",
			ApiKey:              "",
			SystemMessage:       "You are an English tutor who teach English for ESL learner",
			AssistanceMessage:   "Need to give student questions to learn vocabulary. There are 3 types: Definition, Synonym and Application. Here are some examples for each type in JSON format, the example is wrapped by triple backtick (```):\n1. \"Definition\" question: ```{\"questionStem\":\"**cone** means\",    \"questionType\": \"definition\",\"answers\": [\"a shape with a circular base and sides tapering to a point\"],\"distractors\":[\"audacious behavior that you have no right to\",\"a wide scope\",\"a quality that arouses emotions\"]}```\n2. \"Synonym\": question:```{\"questionStem\":\"**diminish** has the same or almost the same meaning as\",\"questionType\":\"synonym\",\"answers\":[\"fall\"],\"distractors\":[\"restore\",\"quiver\",\"glare\"]}```\n3. \"Application\" questions:```[{\"questionStem\":\"Which of the following would most likely **glint**?\",\"questionType\":\"application\",\"answers\":[\"a newly polished car in the sunlight\"],\"distractors\":[\"a newly written manuscript on an editor’s desk\",\"an old sweater in a trunk in the attic\",\"a fine powder used to absorb moisture\"]},{\"questionStem\":\"In which scenario would you most likely use the word **apparently**?\",\"questionType\":\"application\",\"answers\":[\"when you hear a surprising piece of information that you cannot verify\"],\"distractors\":[\"when you are certain about the outcome of an event\",\"when you want to express appreciation for someone's help\",\"when you are describing an experience that happened in the past\"]},{\"questionStem\":\"Which situation best illustrates the meaning of the word **respective**?\",\"questionType\":\"application\",\"answers\":[\"two friends talking about their hobbies: one enjoys painting, and the other enjoys playing guitar\"],\"distractors\":[\"a family enjoying a picnic in the park\",\"a group of students studying together for an exam\",\"a team working on a project in the office\"]}]```.",
			UserMessageTemplate: "Give a single choice question with %s type to help student learn word **%s**. Provide the question in JSON format with keys: `questionStem`, `questionType`, `answers`, `distractors`. Note that, there must be only one correct answer (length of `answers` array is 1), `distractors` must be 3 wrong choices, the word **%s** must appear in field `questionStem` and not appear in 2 fields `answers` and `distractors`.",
			MaxTokens:           256,
			Temperature:         1,
			TopP:                1,
		},
		GeminiClient: GeminiClient{
			Model:       "gemini-1.0-pro",
			ApiKey:      "",
			Temperature: 1,
			TopP:        1,
			MaxTokens:   1024,
			PromptParts: []string{
				"You're an experienced English tutor. Compose exercises for ESL learners of C1 level to reinforce their vocabulary, via multiple-choice questions which have one correct answer and three incorrect choices. Three question types are allowed. A question of DEFINITION type has definitions as choices. A question of SYNONYM type has potential synonyms as choices. A question of APPLICATION type has choices that reflect different scenarios or usage of a concept. The output must be a single object in JSON format like in my given examples. Note that, the correct choice is in \"answers\" field, three incorrect choices are in \"distractors\" field.",
				"word: cone",
				"questionType: definition",
				"output: {\n    \"questionStem\": \"**cone** means\",\n    \"questionType\": \"definition\",\n    \"answers\": [\"a shape with a circular base and sides tapering to a point\"],\n    \"distractors\": [\n      \"audacious behavior that you have no right to\",\n      \"a wide scope\",\n      \"a quality that arouses emotions\"\n    ]\n  }",
				"word: diminish",
				"questionType: synonym",
				"output: {\n    \"questionStem\": \"**“diminish”** has the same or almost the same meaning as\",\n    \"questionType\": \"synonym\",\n    \"answers\": [\"fall\"],\n    \"distractors\": [\"restore\", \"quiver\", \"glare\"],\n  }",
				"word: glint",
				"questionType: application",
				"output: {\n    \"questionStem\": \"Which of the following would most likely **glint**?\",\n    \"questionType\": \"application\",\n    \"answers\": [\"a newly polished car in the sunlight\"],\n    \"distractors\": [\n      \"a newly written manuscript on an editor’s desk\",\n      \"an old sweater in a trunk in the attic\",\n      \"a fine powder used to absorb moisture\"\n    ]\n  }",
				"word: apparently",
				"questionType: ",
				"output: {\n    \"questionStem\": \"In which scenario would you most likely use the word **apparently**?\",\n    \"questionType\": \"application\",\n    \"answers\": [\"when you hear a surprising piece of information that you cannot verify\"],\n    \"distractors\": [\n      \"when you are certain about the outcome of an event\",\n      \"when you want to express appreciation for someone's help\",\n      \"when you are describing an experience that happened in the past\"\n    ]\n  }",
				"word: respective",
				"questionType: application",
				"output: {\n    \"questionStem\": \"Which situation best illustrates the meaning of the word **respective**?\",\n    \"questionType\": \"application\",\n    \"answers\": [\"two friends talking about their hobbies: one enjoys painting, and the other enjoys playing guitar\"],\n    \"distractors\": [\n      \"a family enjoying a picnic in the park\",\n      \"a group of students studying together for an exam\",\n      \"a team working on a project in the office\"\n    ]\n  }",
				"word: proportionately",
				"questionType: application",
				"output: {\n    \"questionStem\": \"Which of the following demonstrates the best use of the word **proportionately**?\",\n    \"questionType\": \"application\",\n    \"answers\": [\"distributing a bonus among employees based on their performance\"],\n    \"distractors\": [\n      \"distributing the tasks equally among the team members\",\n      \"giving a gift to a friend who has helped you greatly\",\n      \"dividing the cost of a meal among the diners evenly\"\n    ]\n  }",
				"word: corresponding",
				"questionType: application",
				"output: {\n    \"questionStem\": \"Which of the following best illustrates the meaning of the word **corresponding**?\",\n    \"questionType\": \"application\",\n    \"answers\": [\"matching the names on the list to the addresses\"],\n    \"distractors\": [\n      \"exchanging gifts with a friend\",\n      \"comparing two different versions of a document\",\n      \"following the instructions in a recipe\"\n    ]\n  }",
				"word: redundant",
				"questionType: definition",
				"output: {\n    \"questionStem\": \"**redundant** means\",\n    \"questionType\": \"definition\",\n    \"answers\": [\"more than is needed, desired, or required\"],\n    \"distractors\": [\"very important\", \"extra help or support that you can get if necessary\n\", \"more than is usual, expected, or than exists already\"]\n  }",
				"word: redundant",
				"questionType: application",
				"output: {\n    \"questionStem\": \"In what scenario would you most likely use the word **redundant**?\",\n    \"questionType\": \"application\",\n    \"answers\": [\n        \"When referring to information or processes that are unnecessary or repetitive\"\n    ],\n    \"distractors\": [\n        \"When describing a situation where resources are insufficient\",\n        \"When discussing a task that requires careful attention to detail\",\n        \"When explaining a concept that is essential to understanding the topic\"\n    ]\n}",
			},
		},
		DataFolder:                "data",
		EvaluationFilePath:        "/home/kiet13/KiT/exercise-generator/data",
		EvaluationIntervalSeconds: 1,
	}
}

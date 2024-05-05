package config

// Config application
type Config struct {
	OpenaiClient OpenaiClient `json:"openai_client"`

	DataFolder                string `json:"data_folder"`
	EvaluationFilePath        string `json:"evaluation_file_path"`
	EvaluationIntervalSeconds int    `json:"evaluation_interval_seconds"`
}

type OpenaiClient struct {
	Host                string  `json:"host"`
	ApiKey              string  `json:"api_key"`
	SystemMessage       string  `json:"system_message"`
	AssistanceMessage   string  `json:"assistance_message"`
	UserMessageTemplate string  `json:"user_message_template"`
	MaxTokens           int     `json:"max_tokens"`
	Temperature         float32 `json:"temperature"`
	TopP                float32 `json:"top_p"`
}

func Load() (*Config, error) {
	// TODO: load config from env variables
	return loadDefaultConfig(), nil
}

// nolint
func loadDefaultConfig() *Config {
	c := &Config{
		OpenaiClient: OpenaiClient{
			Host:                "https://api.openai.com",
			ApiKey:              "sk-proj-MpqloVzjiROlyTsnqJfTT3BlbkFJN52AbvMZHwefwQqQxMXI",
			SystemMessage:       "You are an English tutor who teach English for ESL learner",
			AssistanceMessage:   "Need to give student questions to learn vocabulary. There are 3 types: Definition, Synonym and Application. Here are some examples for each type in JSON format, the example is wrapped by triple backtick (```):\n1. \"Definition\" question: ```{\"questionStem\":\"**cone** means\",    \"questionType\": \"definition\",\"answers\": [\"a shape with a circular base and sides tapering to a point\"],\"distractors\":[\"audacious behavior that you have no right to\",\"a wide scope\",\"a quality that arouses emotions\"]}```\n2. \"Synonym\": question:```{\"questionStem\":\"**diminish** has the same or almost the same meaning as\",\"questionType\":\"synonym\",\"answers\":[\"fall\"],\"distractors\":[\"restore\",\"quiver\",\"glare\"]}```\n3. \"Application\" questions:```[{\"questionStem\":\"Which of the following would most likely **glint**?\",\"questionType\":\"application\",\"answers\":[\"a newly polished car in the sunlight\"],\"distractors\":[\"a newly written manuscript on an editor’s desk\",\"an old sweater in a trunk in the attic\",\"a fine powder used to absorb moisture\"]},{\"questionStem\":\"In which scenario would you most likely use the word **apparently**?\",\"questionType\":\"application\",\"answers\":[\"when you hear a surprising piece of information that you cannot verify\"],\"distractors\":[\"when you are certain about the outcome of an event\",\"when you want to express appreciation for someone's help\",\"when you are describing an experience that happened in the past\"]},{\"questionStem\":\"Which situation best illustrates the meaning of the word **respective**?\",\"questionType\":\"application\",\"answers\":[\"two friends talking about their hobbies: one enjoys painting, and the other enjoys playing guitar\"],\"distractors\":[\"a family enjoying a picnic in the park\",\"a group of students studying together for an exam\",\"a team working on a project in the office\"]}]```.",
			UserMessageTemplate: "Give a single choice question with %s type to help student learn word **%s**. Provide the question in JSON format with keys: `questionStem`, `questionType`, `answers`, `distractors`. Note that, there must be only one correct answer (length of `answers` array is 1), `distractors` must be 3 wrong choices, the word **%s** must appear in field `questionStem` and not appear in 2 fields `answers` and `distractors`.",
			MaxTokens:           256,
			Temperature:         1,
			TopP:                1,
		},
		DataFolder:                "data",
		EvaluationFilePath:        "/home/kiet13/KiT/exercise-generator/data",
		EvaluationIntervalSeconds: 4,
	}
	return c
}

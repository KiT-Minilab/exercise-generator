package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"exercise-generator/config"
	"exercise-generator/internal/api"

	"github.com/urfave/cli/v3"
	"go.uber.org/zap"
)

var (
	cfg    *config.Config
	logger *zap.Logger
)

func main() {
	var err error
	logger, err = zap.NewProduction()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	cfg, err = config.Load()
	if err != nil {
		logger.Fatal("failed to load config", zap.Error(err))
	}

	cmd := &cli.Command{
		Commands: []*cli.Command{
			{
				Name:   "server",
				Usage:  "Generate evaluation template for baseline words in csv format",
				Action: startServer,
			},
			{
				Name:     "migrate",
				Usage:    "Doing database migration",
				Commands: GetMigrationCommands(cfg.MigrationFolder, cfg.PostgreSQL.String()),
			},
			{
				Name:  "generate-evaluation-template",
				Usage: "Generate evaluation template for baseline words in csv format",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "questionType", Value: "", Usage: "specify the question type (definition, synonym, application)"},
					&cli.StringFlag{Name: "provider", Value: "OPENAI", Usage: "specify the AI provider (OPENAI, GOOGLE)"},
				},
				Action: evaluateBaselineWords,
			},
		},
	}

	err = cmd.Run(context.Background(), os.Args)
	if err != nil {
		log.Fatalln("Error:", err)
	}
}

func startServer(ctx context.Context, cmd *cli.Command) error {
	h := api.NewApiHandler(logger, cfg)

	http.HandleFunc("GET /hello", h.HelloWorld)
	http.HandleFunc("GET /api/v1/generation-configs", h.ListGenerationConfig)
	http.HandleFunc("POST /api/v1/generation-configs", h.CreateGenerationConfig)
	http.HandleFunc("PUT /api/v1/generation-configs/{id}", h.UpdateGenerationConfig)

	log.Println("Starting server on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalln("Error starting server:", err)
	}
	return nil
}

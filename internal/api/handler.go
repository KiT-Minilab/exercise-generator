package api

import (
	"exercise-generator/config"
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

type ApiHandler struct {
	log *zap.Logger
	cfg *config.Config
}

func NewApiHandler(logger *zap.Logger, cfg *config.Config) *ApiHandler {
	return &ApiHandler{
		log: logger,
		cfg: cfg,
	}
}

func (h *ApiHandler) HelloWorld(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
}

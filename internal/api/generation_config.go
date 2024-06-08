package api

import (
	"net/http"

	"github.com/spf13/cast"
	"go.uber.org/zap"

	"exercise-generator/internal/model"
)

type Messages struct {
	SystemMessage       string `json:"systemMessage"`
	AssistantMessage    string `json:"assistantMessage"`
	UserMessageTemplate string `json:"userMessageTemplate"`
}

func (h *ApiHandler) CreateGenerationConfig(w http.ResponseWriter, r *http.Request) {
	var genCfg model.GenerationConfig

	err := h.readJSON(w, r, &genCfg)
	if err != nil {
		h.writeJSON(w, http.StatusBadRequest, map[string]interface{}{"message": err.Error()}, nil)
		return
	}

	err = h.service.CreateGenerationConfig(r.Context(), genCfg)
	if err != nil {
		h.log.Error("failed to create generation config", zap.Error(err))
		h.writeJSON(w, http.StatusInternalServerError, map[string]interface{}{"message": "failed to create generation config"}, nil)
		return
	}

	h.writeJSON(w, http.StatusCreated, map[string]interface{}{"message": "success"}, nil)
}

func (h *ApiHandler) UpdateGenerationConfig(w http.ResponseWriter, r *http.Request) {
	var genCfg model.GenerationConfig

	err := h.readJSON(w, r, &genCfg)
	if err != nil {
		h.writeJSON(w, http.StatusBadRequest, map[string]interface{}{"message": err.Error()}, nil)
		return
	}

	genCfg.ID = cast.ToInt((r.PathValue("id")))

	err = h.service.UpdateGenerationConfig(r.Context(), genCfg)
	if err != nil {
		h.log.Error("failed to update generation config", zap.Error(err))
		h.writeJSON(w, http.StatusInternalServerError, map[string]interface{}{"message": "failed to update generation config"}, nil)
		return
	}

	h.writeJSON(w, http.StatusOK, map[string]interface{}{"message": "success"}, nil)
}

func (h *ApiHandler) ListGenerationConfig(w http.ResponseWriter, r *http.Request) {
	genCfgs, err := h.service.ListGenerationConfig(r.Context())
	if err != nil {
		h.log.Error("failed to list generation config", zap.Error(err))
		h.writeJSON(w, http.StatusInternalServerError, map[string]interface{}{"message": "failed to list generation config"}, nil)
		return
	}

	h.writeJSON(w, http.StatusOK, map[string]interface{}{"message": "OK", "generationConfigs": genCfgs}, nil)
}

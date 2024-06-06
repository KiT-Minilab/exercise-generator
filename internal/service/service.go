package service

import (
	"exercise-generator/config"

	"go.uber.org/zap"
)

type Service struct {
	log *zap.Logger
	cfg *config.Config
}

func NewService(logger *zap.Logger, cfg *config.Config) *Service {
	return &Service{
		log: logger,
		cfg: cfg,
	}
}

func (s *Service) HelloWorld() {
}

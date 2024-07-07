package repository

import (
	"exercise-generator/config"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type IRepository interface {
	IGenerationConfigRepository
}

type repositoryImpl struct {
	db  *sqlx.DB
	cfg *config.Config
}

func NewRepository(cfg *config.Config) (IRepository, error) {

	db, err := sqlx.Connect("postgres", cfg.PostgreSQL.String())
	if err != nil {
		return nil, err
	}

	return &repositoryImpl{db: db, cfg: cfg}, nil
}

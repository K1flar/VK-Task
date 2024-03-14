package postgres

import (
	"database/sql"
	"film_library/internal/config"
	"fmt"

	_ "github.com/lib/pq"
)

type Repository struct {
}

func New(cfg *config.DataBase) (*Repository, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &Repository{}, nil
}

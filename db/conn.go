package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mauricioabreu/a-fast-api/config"
)

func NewDB(cfg *config.Config) (*pgxpool.Pool, error) {
	url := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?default_query_exec_mode=simple_protocol",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName)
	pool, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, err
	}

	db, err := pgxpool.NewWithConfig(context.Background(), pool)
	if err != nil {
		return nil, err
	}

	return db, nil
}

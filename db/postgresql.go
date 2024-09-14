package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Storage struct {
	DB *pgxpool.Pool
}

func New(storagePath string) (*Storage, error) {
	db, err := pgxpool.Connect(context.Background(), storagePath)
	if err != nil {
		return nil, fmt.Errorf("could not connect to database: %w", err)
	}
	return &Storage{DB: db}, nil
}

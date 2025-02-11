package repository

import (
	"context"
	"fmt"

	logger "github.com/bllooop/coinshop/pkg/logging"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBname   string
	SSLMode  string
}

const (
	userListTable     = "userlist"
	shopTable         = "shop"
	transactionsTable = "transactions"
	purchaseTable     = "purchases"
)

func NewPostgresDB(cfg Config) (*pgxpool.Pool, error) {
	logger.Log.Info().Msg("Подключение к базе данных")
	db, err := pgxpool.New(context.Background(), fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DBname, cfg.SSLMode))

	if err != nil {
		return nil, err
	}
	return db, nil
}

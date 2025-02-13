package repository

import (
	"fmt"

	logger "github.com/bllooop/coinshop/pkg/logging"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
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

func NewPostgresDB(cfg Config) (*sqlx.DB, error) {
	logger.Log.Info().Msg("Подключение к базе данных")
	db, err := sqlx.Open("postgres", fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DBname, cfg.SSLMode))
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DBname, cfg.SSLMode)
	logger.Log.Info().Msgf("Connecting to database at: %s", connStr)
	if err != nil {
		return nil, err
	}
	return db, nil
}

package repository

import (
	"database/sql"
	"fmt"

	logger "github.com/bllooop/coinshop/pkg/logging"
	_ "github.com/lib/pq"
	"github.com/pressly/goose"
)

func RunMigrate(cfg Config, migratePath string) error {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DBname, cfg.SSLMode)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return err
	}
	defer db.Close()
	logger.Log.Info().Msg("Applying migrations / Применение миграций")
	err = goose.SetDialect("postgres")
	if err != nil {
		return err
	}
	err = goose.Up(db, migratePath)
	if err != nil {
		return err
	}
	logger.Log.Info().Msg("Migrations applied successfully! / Миграция прошла успешно!")
	return nil
}

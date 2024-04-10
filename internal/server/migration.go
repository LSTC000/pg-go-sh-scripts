package server

import (
	"database/sql"
	"fmt"
	"pg-sh-scripts/internal/common"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

const (
	migrationDir       = "migration"
	pgMigrationDialect = "postgres"
)

func setMigration(pool *pgxpool.Pool) error {
	if err := goose.SetDialect(pgMigrationDialect); err != nil {
		return err
	}

	db := stdlib.OpenDBFromPool(pool)
	defer func(db *sql.DB) {
		if err := db.Close(); err != nil {
			logger := common.GetLogger()
			logger.Error(fmt.Sprintf("Close migration db error: %v", err))
		}
	}(db)

	if err := goose.Up(db, migrationDir); err != nil {
		return err
	}

	return nil
}

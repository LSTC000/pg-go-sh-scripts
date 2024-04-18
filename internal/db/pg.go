package db

import (
	"context"
	"pg-sh-scripts/internal/config"
	"pg-sh-scripts/pkg/client/postgres"
	"sync"
)

var (
	pgInstance postgres.IClient
	pgConnErr  error
	pgOnce     sync.Once
)

func GetPgClient() (postgres.IClient, error) {
	pgOnce.Do(func() {
		cfg := config.GetConfig()

		connConfig := postgres.ConnConfig{
			Database:          cfg.Postgres.Database,
			Username:          cfg.Postgres.Username,
			Password:          cfg.Postgres.Password,
			Host:              cfg.Postgres.Host,
			Port:              cfg.Postgres.Port,
			RetryCount:        cfg.Postgres.RetryCount,
			RetrySleepSeconds: cfg.Postgres.RetrySleepSeconds,
		}

		client, err := postgres.GetClient(context.Background(), &connConfig)
		if err != nil {
			pgConnErr = err
			return
		}
		pgInstance = client
	})

	if pgConnErr != nil {
		return nil, pgConnErr
	}

	return pgInstance, nil
}

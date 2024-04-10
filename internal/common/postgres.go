package common

import (
	"context"
	"pg-sh-scripts/internal/config"
	"pg-sh-scripts/pkg/client/postgres"
	"sync"
)

func GetPgClient() (postgres.IClient, error) {
	var (
		pgClient postgres.IClient
		connErr  error
		once     sync.Once
	)

	cfg := config.GetConfig()

	once.Do(func() {
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
			connErr = err
			return
		}
		pgClient = client
	})

	if connErr != nil {
		return nil, connErr
	}

	return pgClient, nil
}

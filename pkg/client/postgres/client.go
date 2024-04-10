package postgres

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type (
	IClient interface {
		GetDB() *pgxpool.Pool
		Close()
	}

	Client struct {
		db *pgxpool.Pool
	}

	ConnConfig struct {
		Database          string
		Username          string
		Password          string
		Host              string
		Port              string
		RetryCount        int
		RetrySleepSeconds time.Duration
	}
)

func getConnString(connConfig *ConnConfig) string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		connConfig.Username,
		connConfig.Password,
		connConfig.Host,
		connConfig.Port,
		connConfig.Database,
	)
}

func getPool(ctx context.Context, connConfig *ConnConfig) (*pgxpool.Pool, error) {
	var (
		poolInstance *pgxpool.Pool
		connErr      error
	)

	db, err := pgxpool.New(ctx, getConnString(connConfig))
	if err != nil {
		return nil, err
	}

	connRetryCounter := 0
	for connRetryCounter < connConfig.RetryCount {
		log.Printf("Connecting to Postgres: try %d/%d", connRetryCounter, connConfig.RetryCount)
		if err := db.Ping(ctx); err != nil {
			connErr = err
			connRetryCounter++
			time.Sleep(connConfig.RetrySleepSeconds)
		} else {
			log.Print("Successful connection to Postgres")
			connErr = nil
			poolInstance = db
			connRetryCounter = connConfig.RetryCount
		}
	}

	return poolInstance, connErr
}

func (pg *Client) GetDB() *pgxpool.Pool {
	return pg.db
}

func (pg *Client) Close() {
	pg.db.Close()
}

func GetClient(ctx context.Context, connConfig *ConnConfig) (IClient, error) {
	db, err := getPool(ctx, connConfig)
	if err != nil {
		return nil, err
	}
	return &Client{db}, nil
}

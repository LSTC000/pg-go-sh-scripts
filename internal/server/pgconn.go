package server

import (
	"pg-sh-scripts/internal/db"
	"pg-sh-scripts/pkg/client/postgres"
)

func setPgConn() (postgres.IClient, error) {
	pgClient, err := db.GetPgClient()
	if err != nil {
		return nil, err
	}
	return pgClient, nil
}

func closePgConn() error {
	pgClient, err := db.GetPgClient()
	if err != nil {
		return err
	}
	pgClient.Close()
	return nil
}

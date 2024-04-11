package bash

import (
	"context"
	"errors"
	"fmt"
	"pg-sh-scripts/internal/common"
	"pg-sh-scripts/pkg/logging"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	uuid "github.com/satori/go.uuid"
)

type PgRepository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

func (p PgRepository) GetOneById(ctx context.Context, id uuid.UUID) (*Bash, error) {
	bash := Bash{}

	p.logger.Debug(fmt.Sprintf("Start getting bash by id: %s", id))
	q := `
		SELECT
			id, title, body, created_at
		FROM
		    scripts.bash
		WHERE 
			id = $1
	`

	row := p.db.QueryRow(ctx, q, id)
	if err := row.Scan(&bash.Id, &bash.Title, &bash.Body, &bash.CreatedAt); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			p.logger.Error(fmt.Sprintf("Getting bash Error: %s, Detail: %s, Where: %s", pgErr.Message, pgErr.Detail, pgErr.Where))
		}
		return nil, err
	}
	p.logger.Debug(fmt.Sprintf("Finish getting bash by id: %s", id))

	return &bash, nil
}

func (p PgRepository) GetAll(ctx context.Context) ([]*Bash, error) {
	bashList := make([]*Bash, 0)

	p.logger.Debug("Start getting bash list")
	q := `
		SELECT
			id, title, body, created_at
		FROM
		    scripts.bash
	`

	rows, err := p.db.Query(ctx, q)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			p.logger.Error(fmt.Sprintf("Getting bash list Error: %s, Detail: %s, Where: %s", pgErr.Message, pgErr.Detail, pgErr.Where))
		}
		return nil, err
	}

	for rows.Next() {
		bash := Bash{}
		if err := rows.Scan(&bash.Id, &bash.Title, &bash.Body, &bash.CreatedAt); err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) {
				p.logger.Error(fmt.Sprintf("Getting bash Error: %s, Detail: %s, Where: %s", pgErr.Message, pgErr.Detail, pgErr.Where))
			}
			return nil, err
		}
		bashList = append(bashList, &bash)
	}
	p.logger.Debug("Finish getting bash list")

	return bashList, nil
}

func (p PgRepository) Create(ctx context.Context, createBash CreateBashDTO) (*Bash, error) {
	bash := Bash{}

	p.logger.Debug(fmt.Sprintf("Start creating bash with title: %s", createBash.Title))
	stmt := `
		INSERT INTO scripts.bash
			(title, body)
		VALUES 
			($1, $2)
		RETURNING id, title, body, created_at
	`

	row := p.db.QueryRow(ctx, stmt, createBash.Title, createBash.Body)
	if err := row.Scan(&bash.Id, &bash.Title, &bash.Body, &bash.CreatedAt); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			p.logger.Error(fmt.Sprintf("Creating bash Error: %s, Detail: %s, Where: %s", pgErr.Message, pgErr.Detail, pgErr.Where))
		}
		return nil, err
	}
	p.logger.Debug(fmt.Sprintf("Finish creating bash with title: %s", createBash.Title))

	return &bash, nil
}

func GetPgRepository() IRepository {
	logger := common.GetLogger()
	pg, err := common.GetPgClient()
	if err != nil {
		logger.Error(fmt.Sprintf("Getting postgres client Error: %s", err))
		panic(err)
	}
	return &PgRepository{
		db:     pg.GetDB(),
		logger: logger,
	}
}

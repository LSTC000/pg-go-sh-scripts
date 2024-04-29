package repo

import (
	"context"
	"errors"
	"fmt"
	"pg-sh-scripts/internal/db"
	"pg-sh-scripts/internal/dto"
	"pg-sh-scripts/internal/log"
	"pg-sh-scripts/internal/model"
	"pg-sh-scripts/internal/type/alias"
	"pg-sh-scripts/pkg/logging"
	"pg-sh-scripts/pkg/sql/pagination"

	"github.com/georgysavva/scany/v2/pgxscan"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	uuid "github.com/satori/go.uuid"
)

type PgBashRepository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

func (p PgBashRepository) GetOneById(ctx context.Context, id uuid.UUID) (*model.Bash, error) {
	bash := &model.Bash{}

	p.logger.Debug(fmt.Sprintf("Start getting bash by id: %v", id))
	q := `
		SELECT
			id, title, body, created_at
		FROM
		    scripts.bash
		WHERE 
			id = $1
	`

	if err := pgxscan.Get(ctx, p.db, bash, q, id); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			p.logger.Error(
				fmt.Sprintf(
					"Getting bash by id: %v Error: %s, Detail: %s, Where: %s",
					id,
					pgErr.Message,
					pgErr.Detail,
					pgErr.Where,
				),
			)
		} else {
			p.logger.Error(fmt.Sprintf("Getting bash by id: %v Error: %s", id, err))
		}
		return bash, err
	}
	p.logger.Debug(fmt.Sprintf("Finish getting bash by id: %v", id))

	return bash, nil
}

func (p PgBashRepository) GetPaginationPage(
	ctx context.Context,
	paginationParams pagination.LimitOffsetParams,
) (alias.BashLimitOffsetPage, error) {
	var bashPaginationPage alias.BashLimitOffsetPage

	p.logger.Debug("Start getting bash pagination page")
	q := `
		SELECT
			id, title, body, created_at
		FROM
		    scripts.bash
	`

	bashPaginationPage, err := pagination.Paginate[*model.Bash](ctx, p.db, q, paginationParams)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			p.logger.Error(
				fmt.Sprintf(
					"Getting bash pagination page Error: %s, Detail: %s, Where: %s",
					pgErr.Message,
					pgErr.Detail,
					pgErr.Where,
				),
			)
		} else {
			p.logger.Error(fmt.Sprintf("Getting bash pagination page Error: %s", err))
		}
		return bashPaginationPage, err
	}
	p.logger.Debug("Finish getting bash pagination page")

	return bashPaginationPage, nil
}

func (p PgBashRepository) Create(ctx context.Context, dto dto.CreateBash) (*model.Bash, error) {
	bash := &model.Bash{}

	p.logger.Debug(fmt.Sprintf("Start creating bash with title: %s", dto.Title))
	stmt := `
		INSERT INTO scripts.bash
			(title, body)
		VALUES 
			($1, $2)
		RETURNING id, title, body, created_at
	`

	if err := pgxscan.Get(ctx, p.db, bash, stmt, dto.Title, dto.Body); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			p.logger.Error(
				fmt.Sprintf(
					"Creating bash with title: %s Error: %s, Detail: %s, Where: %s",
					dto.Title,
					pgErr.Message,
					pgErr.Detail,
					pgErr.Where,
				),
			)
		} else {
			p.logger.Error(fmt.Sprintf("Creating bash with title: %s Error: %s", dto.Title, err))
		}
		return bash, err
	}
	p.logger.Debug(fmt.Sprintf("Finish creating bash with title: %s", dto.Title))

	return bash, nil
}

func (p PgBashRepository) RemoveById(ctx context.Context, id uuid.UUID) (*model.Bash, error) {
	bash := &model.Bash{}

	p.logger.Debug(fmt.Sprintf("Start removing bash by id: %v", id))
	stmt := `
		DELETE FROM 
		    scripts.bash
		WHERE 
			id = $1
		RETURNING id, title, body, created_at
	`

	if err := pgxscan.Get(ctx, p.db, bash, stmt, id); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			p.logger.Error(
				fmt.Sprintf(
					"Removing bash by id: %v Error: %s, Detail: %s, Where: %s",
					id,
					pgErr.Message,
					pgErr.Detail,
					pgErr.Where,
				),
			)
		} else {
			p.logger.Error(fmt.Sprintf("Removing bash by id: %v Error: %s", id, err))
		}
		return bash, err
	}
	p.logger.Debug(fmt.Sprintf("Finish removing bash by id: %v", id))

	return bash, nil
}

func GetPgBashRepository() IBashRepository {
	logger := log.GetLogger()
	pg, err := db.GetPgClient()
	if err != nil {
		logger.Error(fmt.Sprintf("Getting postgres client Error: %s", err))
		panic(err)
	}
	return &PgBashRepository{
		db:     pg.GetDB(),
		logger: logger,
	}
}

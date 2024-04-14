package repo

import (
	"context"
	"errors"
	"fmt"
	"pg-sh-scripts/internal/db"
	"pg-sh-scripts/internal/dto"
	"pg-sh-scripts/internal/log"
	"pg-sh-scripts/internal/model"
	"pg-sh-scripts/internal/schema"
	"pg-sh-scripts/pkg/logging"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	uuid "github.com/satori/go.uuid"
)

type PgBashRepository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

func (p PgBashRepository) GetOneById(ctx context.Context, id uuid.UUID) (*model.Bash, error) {
	bash := model.Bash{}

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

func (p PgBashRepository) GetPaginationPage(ctx context.Context, paginationParams schema.PaginationParams) (schema.PaginationPage[*model.Bash], error) {
	limit := paginationParams.Limit
	offset := paginationParams.Offset

	bashList := make([]*model.Bash, 0, limit)
	bashPaginationPage := schema.PaginationPage[*model.Bash]{
		Limit:  limit,
		Offset: offset,
	}

	p.logger.Debug("Start getting bash list")
	qItems := `
		SELECT
			id, title, body, created_at
		FROM
		    scripts.bash
		OFFSET $1
		LIMIT $2
	`
	qTotal := `
		SELECT
			COUNT(*) AS total
		FROM
		    scripts.bash
	`

	rows, err := p.db.Query(ctx, qItems, offset, limit)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			p.logger.Error(fmt.Sprintf("Getting bash items Error: %s, Detail: %s, Where: %s", pgErr.Message, pgErr.Detail, pgErr.Where))
		}
		return bashPaginationPage, err
	}
	for rows.Next() {
		bash := model.Bash{}
		if err := rows.Scan(&bash.Id, &bash.Title, &bash.Body, &bash.CreatedAt); err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) {
				p.logger.Error(fmt.Sprintf("Getting bash Error: %s, Detail: %s, Where: %s", pgErr.Message, pgErr.Detail, pgErr.Where))
			}
			return bashPaginationPage, err
		}
		bashList = append(bashList, &bash)
	}
	bashPaginationPage.Items = bashList

	row := p.db.QueryRow(ctx, qTotal)
	if err := row.Scan(&bashPaginationPage.Total); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			p.logger.Error(fmt.Sprintf("Getting bash total Error: %s, Detail: %s, Where: %s", pgErr.Message, pgErr.Detail, pgErr.Where))
		}
		return bashPaginationPage, err
	}
	p.logger.Debug("Finish getting bash list")

	return bashPaginationPage, nil
}

func (p PgBashRepository) Create(ctx context.Context, dto dto.CreateBashDTO) (*model.Bash, error) {
	bash := model.Bash{}

	p.logger.Debug(fmt.Sprintf("Start creating bash with title: %s", dto.Title))
	stmt := `
		INSERT INTO scripts.bash
			(title, body)
		VALUES 
			($1, $2)
		RETURNING id, title, body, created_at
	`

	row := p.db.QueryRow(ctx, stmt, dto.Title, dto.Body)
	if err := row.Scan(&bash.Id, &bash.Title, &bash.Body, &bash.CreatedAt); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			p.logger.Error(fmt.Sprintf("Creating bash Error: %s, Detail: %s, Where: %s", pgErr.Message, pgErr.Detail, pgErr.Where))
		}
		return nil, err
	}
	p.logger.Debug(fmt.Sprintf("Finish creating bash with title: %s", dto.Title))

	return &bash, nil
}

func (p PgBashRepository) RemoveById(ctx context.Context, id uuid.UUID) (*model.Bash, error) {
	bash := model.Bash{}

	p.logger.Debug(fmt.Sprintf("Start removing bash with id: %s", id))
	stmt := `
		DELETE FROM 
		    scripts.bash
		WHERE 
			id = $1
		RETURNING id, title, body, created_at
	`

	row := p.db.QueryRow(ctx, stmt, id)
	if err := row.Scan(&bash.Id, &bash.Title, &bash.Body, &bash.CreatedAt); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			p.logger.Error(fmt.Sprintf("Removing bash Error: %s, Detail: %s, Where: %s", pgErr.Message, pgErr.Detail, pgErr.Where))
		}
		return nil, err
	}
	p.logger.Debug(fmt.Sprintf("Finish removing bash with id: %s", id))

	return &bash, nil
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

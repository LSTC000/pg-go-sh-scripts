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

	uuid "github.com/satori/go.uuid"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PgBashLogRepository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

func (p PgBashLogRepository) GetPaginationPageByBashId(ctx context.Context, bashId uuid.UUID, paginationParams schema.PaginationParams) (schema.PaginationPage[*model.BashLog], error) {
	limit := paginationParams.Limit
	offset := paginationParams.Offset

	bashLogList := make([]*model.BashLog, 0, limit)
	bashLogPaginationPage := schema.PaginationPage[*model.BashLog]{
		Limit:  limit,
		Offset: offset,
	}

	p.logger.Debug(fmt.Sprintf("Start getting bash log list by bash id: %v", bashId))
	qItems := `
		SELECT
			id, bash_id, body, is_error, created_at
		FROM
		    scripts.bash_log
		WHERE 
		    bash_id = $1
		OFFSET $2
		LIMIT $3
	`
	qTotal := `
		SELECT
			COUNT(*) AS total
		FROM
		    scripts.bash_log
		WHERE 
		    bash_id = $1
	`

	rows, err := p.db.Query(ctx, qItems, bashId, offset, limit)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			p.logger.Error(fmt.Sprintf("Getting bash log items Error: %s, Detail: %s, Where: %s", pgErr.Message, pgErr.Detail, pgErr.Where))
		}
		return bashLogPaginationPage, err
	}

	for rows.Next() {
		bashLog := model.BashLog{}
		if err := rows.Scan(&bashLog.Id, &bashLog.BashId, &bashLog.Body, &bashLog.IsError, &bashLog.CreatedAt); err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) {
				p.logger.Error(fmt.Sprintf("Getting bash log Error: %s, Detail: %s, Where: %s", pgErr.Message, pgErr.Detail, pgErr.Where))
			}
			return bashLogPaginationPage, err
		}
		bashLogList = append(bashLogList, &bashLog)
	}
	bashLogPaginationPage.Items = bashLogList

	row := p.db.QueryRow(ctx, qTotal, bashId)
	if err := row.Scan(&bashLogPaginationPage.Total); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			p.logger.Error(fmt.Sprintf("Getting bash total Error: %s, Detail: %s, Where: %s", pgErr.Message, pgErr.Detail, pgErr.Where))
		}
		return bashLogPaginationPage, err
	}
	p.logger.Debug(fmt.Sprintf("Finish getting bash log list by bash id: %v", bashId))

	return bashLogPaginationPage, nil
}

func (p PgBashLogRepository) Create(ctx context.Context, dto dto.CreateBashLogDTO) (*model.BashLog, error) {
	bashLog := model.BashLog{}

	p.logger.Debug(fmt.Sprintf("Start creating bash log for bash with id: %v", dto.BashId))
	stmt := `
		INSERT INTO scripts.bash_log
			(bash_id, body, is_error)
		VALUES 
			($1, $2, $3)
		RETURNING id, bash_id, body, is_error, created_at
	`

	row := p.db.QueryRow(ctx, stmt, dto.BashId, dto.Body, dto.IsError)
	if err := row.Scan(&bashLog.Id, &bashLog.BashId, &bashLog.Body, &bashLog.IsError, &bashLog.CreatedAt); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			p.logger.Error(fmt.Sprintf("Creating bash log Error: %s, Detail: %s, Where: %s", pgErr.Message, pgErr.Detail, pgErr.Where))
		}
		return nil, err
	}
	p.logger.Debug(fmt.Sprintf("Finish creating bash log for bash with id: %v", dto.BashId))

	return &bashLog, nil
}

func GetPgBashLogRepository() IBashLogRepository {
	logger := log.GetLogger()
	pg, err := db.GetPgClient()
	if err != nil {
		logger.Error(fmt.Sprintf("Getting postgres client Error: %s", err))
		panic(err)
	}
	return &PgBashLogRepository{
		db:     pg.GetDB(),
		logger: logger,
	}
}

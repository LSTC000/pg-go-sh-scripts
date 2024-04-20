package repo

import (
	"context"
	"errors"
	"fmt"
	"pg-sh-scripts/internal/db"
	"pg-sh-scripts/internal/dto"
	"pg-sh-scripts/internal/log"
	"pg-sh-scripts/internal/model"
	"pg-sh-scripts/pkg/logging"
	"pg-sh-scripts/pkg/sql/pagination"

	uuid "github.com/satori/go.uuid"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PgBashLogRepository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

func (p PgBashLogRepository) GetPaginationPageByBashId(ctx context.Context, bashId uuid.UUID, paginationParams pagination.LimitOffsetParams) (pagination.LimitOffsetPage[*model.BashLog], error) {
	var bashLogPaginationPage pagination.LimitOffsetPage[*model.BashLog]

	p.logger.Debug(fmt.Sprintf("Start getting bash log pagination page by bash id: %v", bashId))
	q := `
		SELECT
			id, bash_id, body, is_error, created_at
		FROM
		    scripts.bash_log
		WHERE 
		    bash_id = $1
	`

	bashLogPaginationPage, err := pagination.Paginate[*model.BashLog](ctx, p.db, q, paginationParams, bashId)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			p.logger.Error(fmt.Sprintf("Getting bash log pagination page by bash id: %v Error: %s, Detail: %s, Where: %s", bashId, pgErr.Message, pgErr.Detail, pgErr.Where))
		} else {
			p.logger.Error(fmt.Sprintf("Getting bash log pagination page by bash id: %v Error: %s", bashId, err))
		}
		return bashLogPaginationPage, err
	}
	p.logger.Debug(fmt.Sprintf("Finish getting bash log pagination page by bash id: %v", bashId))

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

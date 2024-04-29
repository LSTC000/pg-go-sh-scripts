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

	uuid "github.com/satori/go.uuid"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PgBashLogRepository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

func (p PgBashLogRepository) GetPaginationPageByBashId(
	ctx context.Context,
	bashId uuid.UUID,
	paginationParams pagination.LimitOffsetParams,
) (alias.BashLogLimitOffsetPage, error) {
	var bashLogPaginationPage alias.BashLogLimitOffsetPage

	p.logger.Debug(fmt.Sprintf("Start getting bash log pagination page by bash id: %v", bashId))
	q := `
		SELECT
			id, bash_id, body, is_error, created_at
		FROM
		    scripts.bash_log
		WHERE 
		    bash_id = $1
	`

	bashLogPaginationPage, err := pagination.Paginate[*model.BashLog](
		ctx,
		p.db,
		q,
		paginationParams,
		bashId,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			p.logger.Error(
				fmt.Sprintf(
					"Getting bash log pagination page by bash id: %v Error: %s, Detail: %s, Where: %s",
					bashId,
					pgErr.Message,
					pgErr.Detail,
					pgErr.Where,
				),
			)
		} else {
			p.logger.Error(fmt.Sprintf("Getting bash log pagination page by bash id: %v Error: %s", bashId, err))
		}
		return bashLogPaginationPage, err
	}
	p.logger.Debug(fmt.Sprintf("Finish getting bash log pagination page by bash id: %v", bashId))

	return bashLogPaginationPage, nil
}

func (p PgBashLogRepository) Create(
	ctx context.Context,
	dto dto.CreateBashLog,
) (*model.BashLog, error) {
	bashLog := &model.BashLog{}

	p.logger.Debug(fmt.Sprintf("Start creating bash log by bash id: %v", dto.BashId))
	stmt := `
		INSERT INTO scripts.bash_log
			(bash_id, body, is_error)
		VALUES 
			($1, $2, $3)
		RETURNING id, bash_id, body, is_error, created_at
	`

	if err := pgxscan.Get(ctx, p.db, bashLog, stmt, dto.BashId, dto.Body, dto.IsError); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			p.logger.Error(
				fmt.Sprintf(
					"Creating bash log by bash id: %v Error: %s, Detail: %s, Where: %s",
					dto.BashId,
					pgErr.Message,
					pgErr.Detail,
					pgErr.Where,
				),
			)
		} else {
			p.logger.Error(fmt.Sprintf("Creating bash log by bash id: %v Error: %s", dto.BashId, err))
		}
		return bashLog, err
	}
	p.logger.Debug(fmt.Sprintf("Finish creating bash log by bash id: %v", dto.BashId))

	return bashLog, nil
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

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

	uuid "github.com/satori/go.uuid"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PgBashLogRepository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

func (p PgBashLogRepository) GetAllByBashId(ctx context.Context, bashId uuid.UUID) ([]*model.BashLog, error) {
	bashLogList := make([]*model.BashLog, 0)

	p.logger.Debug(fmt.Sprintf("Start getting bash log list by bash id: %v", bashId))
	q := `
		SELECT
			id, bash_id, body, is_error, created_at
		FROM
		    scripts.bash_log
		WHERE 
		    bash_id = $1
	`

	rows, err := p.db.Query(ctx, q, bashId)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			p.logger.Error(fmt.Sprintf("Getting bash log list Error: %s, Detail: %s, Where: %s", pgErr.Message, pgErr.Detail, pgErr.Where))
		}
		return nil, err
	}

	for rows.Next() {
		bashLog := model.BashLog{}
		if err := rows.Scan(&bashLog.Id, &bashLog.BashId, &bashLog.Body, &bashLog.IsError, &bashLog.CreatedAt); err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) {
				p.logger.Error(fmt.Sprintf("Getting bash log Error: %s, Detail: %s, Where: %s", pgErr.Message, pgErr.Detail, pgErr.Where))
			}
			return nil, err
		}
		bashLogList = append(bashLogList, &bashLog)
	}
	p.logger.Debug(fmt.Sprintf("Finish getting bash log list by bash id: %v", bashId))

	return bashLogList, nil
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
	if err := row.Scan(&bashLog.Id, &bashLog.BashId, &bashLog.Body, &bashLog.CreatedAt); err != nil {
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

package bashlog

import (
	"context"
	"errors"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"pg-sh-scripts/internal/common"
	"pg-sh-scripts/pkg/logging"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PgRepository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

func (p PgRepository) GetAllByBashID(ctx context.Context, bashID uuid.UUID) ([]*BashLog, error) {
	bashLogList := make([]*BashLog, 0)

	p.logger.Debug(fmt.Sprintf("Start getting bash log list by bash id: %v", bashID))
	q := `
		SELECT
			id, bash_id, body, created_at
		FROM
		    scripts.bash_log
		WHERE 
		    bash_id = $1
	`

	rows, err := p.db.Query(ctx, q, bashID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			p.logger.Error(fmt.Sprintf("Getting bash log list Error: %s, Detail: %s, Where: %s", pgErr.Message, pgErr.Detail, pgErr.Where))
		}
		return nil, err
	}

	for rows.Next() {
		bashLog := BashLog{}
		if err := rows.Scan(&bashLog.ID, &bashLog.BashID, &bashLog.Body, &bashLog.CreatedAt); err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) {
				p.logger.Error(fmt.Sprintf("Getting bash log Error: %s, Detail: %s, Where: %s", pgErr.Message, pgErr.Detail, pgErr.Where))
			}
			return nil, err
		}
		bashLogList = append(bashLogList, &bashLog)
	}
	p.logger.Debug(fmt.Sprintf("Finish getting bash log list by bash id: %v", bashID))

	return bashLogList, nil
}

func (p PgRepository) Create(ctx context.Context, createBashLog CreateBashLogDTO) (*BashLog, error) {
	bashLog := BashLog{}

	p.logger.Debug(fmt.Sprintf("Start creating bash log for bash with id: %v", createBashLog.BashID))
	stmt := `
		INSERT INTO scripts.bash_log
			(bash_id, body)
		VALUES 
			($1, $2)
		RETURNING id, bash_id, body, created_at
	`

	row := p.db.QueryRow(ctx, stmt, createBashLog.BashID, createBashLog.Body)
	if err := row.Scan(&bashLog.ID, &bashLog.BashID, &bashLog.Body, &bashLog.CreatedAt); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			p.logger.Error(fmt.Sprintf("Creating bash log Error: %s, Detail: %s, Where: %s", pgErr.Message, pgErr.Detail, pgErr.Where))
		}
		return nil, err
	}
	p.logger.Debug(fmt.Sprintf("Finish creating bash log for bash with id: %v", createBashLog.BashID))

	return &bashLog, nil
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

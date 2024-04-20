package pagination

import (
	"context"
	"fmt"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Paginate[T any](ctx context.Context, db *pgxpool.Pool, query string, params LimitOffsetParams, args ...any) (LimitOffsetPage[T], error) {
	argsCount := len(args)

	itemsArgs := make([]any, 0, argsCount+2)
	totalArgs := make([]any, 0, argsCount)

	offsetArgNumber := argsCount + 1
	limitArgNumber := offsetArgNumber + 1

	if argsCount > 0 {
		for _, arg := range args {
			itemsArgs = append(itemsArgs, arg)
			totalArgs = append(totalArgs, arg)
		}
	}
	itemsArgs = append(itemsArgs, params.Offset, params.Limit)

	items := make([]T, 0, params.Limit)
	page := LimitOffsetPage[T]{
		Limit:  params.Limit,
		Offset: params.Offset,
	}

	qItems := fmt.Sprintf("%s OFFSET $%d LIMIT $%d", query, offsetArgNumber, limitArgNumber)
	qTotal := fmt.Sprintf("SELECT COUNT(*) AS total FROM (%s) AS q1", query)

	if err := pgxscan.Select(ctx, db, &items, qItems, itemsArgs...); err != nil {
		return page, err
	}
	page.Items = items

	if err := pgxscan.Get(ctx, db, &page.Total, qTotal, totalArgs...); err != nil {
		return page, err
	}

	return page, nil
}

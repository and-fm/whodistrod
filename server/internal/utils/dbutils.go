package utils

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func DoesExactlyOneRowExist(pgConn *pgxpool.Pool, ctx context.Context, sql string, args ...any) (bool, error) {
	rows, err := pgConn.Query(ctx, sql, args...)

	if err != nil {
		return false, err
	}

	_, err = pgx.CollectExactlyOneRow(rows, pgx.RowTo[int])

	if err != nil {
		return false, nil
	}

	return true, nil
}

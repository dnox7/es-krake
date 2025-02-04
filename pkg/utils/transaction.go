package utils

import (
	"context"
	"database/sql"
	"pech/es-krake/pkg/log"

	"github.com/jmoiron/sqlx"
)

func Transaction(
	ctx context.Context,
	l *log.Logger,
	db *sqlx.DB,
	opts *sql.TxOptions,
	callback func(tx *sqlx.Tx) error,
) error {
	tx, err := db.BeginTxx(ctx, opts)
	if err != nil {
		l.ErrorContext(ctx, "Failed to begin a transaction", "detail", err)
		return err
	}

	committed := false
	defer (func() {
		if !committed {
			err := tx.Rollback()
			if err != nil {
				l.ErrorContext(ctx, "Failed to rollback transaction", "detail", err)
			}
		}
	})()

	if err := callback(tx); err != nil {
		l.ErrorContext(ctx, "An error occurred while executing the callback", "detail", err)
		return err
	}

	if err := tx.Commit(); err != nil {
		l.ErrorContext(ctx, "Failed to commit the transaction", "detail", err)
		return err
	}

	committed = true
	return nil
}

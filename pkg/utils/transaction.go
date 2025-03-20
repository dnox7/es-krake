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
		l.Error(ctx, "Failed to begin a transaction")
		return err
	}

	committed := false
	defer (func() {
		if !committed {
			err := tx.Rollback()
			if err != nil {
				l.Error(ctx, "Failed to rollback transaction")
			}
		}
	})()

	if err := callback(tx); err != nil {
		l.Error(ctx, "An error occurred while executing the callback")
		return err
	}

	if err := tx.Commit(); err != nil {
		l.Error(ctx, "Failed to commit the transaction")
		return err
	}

	committed = true
	return nil
}

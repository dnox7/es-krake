package utils

import (
	"context"
	"database/sql"
	"pech/es-krake/pkg/log"

	"gorm.io/gorm"
)

func GormTransaction(
	ctx context.Context,
	l *log.Logger,
	db *gorm.DB,
	opts *sql.TxOptions,
	callback func(tx *gorm.DB) error,
) error {
	var err error
	tx := db.WithContext(ctx).Begin(opts)
	if err = tx.Error; err != nil {
		l.Error(ctx, "Failed to begin a transaction")
		return err
	}

	committed := false
	defer (func() {
		if !committed {
			err = tx.Rollback().Error
			if err != nil {
				l.Error(ctx, "Failed to rollback transaction")
			}
		}
	})()

	if err = callback(tx); err != nil {
		l.Error(ctx, "An error occurred while executing the callback")
		return err
	}

	if err = tx.Commit().Error; err != nil {
		l.Error(ctx, "Failed to commit the transaction")
		return err
	}

	committed = true
	return nil
}

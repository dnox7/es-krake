package utils

import (
	"context"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func TestTransaction(t *testing.T) {
	newDB := func() (*gorm.DB, sqlmock.Sqlmock, error) {
		db, dbMock, err := sqlmock.New()
		if err != nil {
			return nil, nil, err
		}

		gormConf := &gorm.Config{
			DisableAutomaticPing: true,
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			},
		}
		gdb, err := gorm.Open(
			postgres.New(postgres.Config{
				Conn: db,
			}),
			gormConf,
		)
		if err != nil {
			return nil, nil, err
		}
		return gdb, dbMock, nil
	}

	t.Run("normal case", func(t *testing.T) {
		db, dbMock, err := newDB()
		if err != nil {
			t.Fatal()
		}
		dbMock.ExpectBegin()
		dbMock.ExpectCommit()

		assert.Nil(t, Transaction(context.Background(), db, func(tx *gorm.DB) error {
			assert.NotNil(t, tx)
			return nil
		}))
		assert.Nil(t, dbMock.ExpectationsWereMet())
	})

	t.Run("error case", func(t *testing.T) {
		db, dbMock, err := newDB()
		if err != nil {
			t.Fatal()
		}

		dbMock.ExpectBegin()
		dbMock.ExpectCommit()

		assert.NotNil(t, Transaction(context.Background(), db, func(tx *gorm.DB) error {
			assert.NotNil(t, tx)
			return fmt.Errorf("error occurred")
		}))
		assert.Nil(t, dbMock.ExpectationsWereMet())
	})
}

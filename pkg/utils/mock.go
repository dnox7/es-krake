package utils

import (
	"database/sql/driver"

	"github.com/DATA-DOG/go-sqlmock"
)

func CreateSingleMockRow(rows map[string]driver.Value) (rowOutput *sqlmock.Rows) {
	var columns []string
	var values []driver.Value
	for v, i := range rows {
		columns = append(columns, v)
		values = append(values, i)
	}
	rowOutput = sqlmock.NewRows(columns).AddRow(values...)
	return
}

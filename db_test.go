package sqlh_test

import (
	"context"
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
)

var db *sql.DB
var ctx = context.Background()

func init() {
	// create mock database for examples
	var mock sqlmock.Sqlmock
	var err error
	db, mock, err = sqlmock.New()
	if err != nil {
		panic(err)
	}
	mock.ExpectQuery("SELECT location FROM binary_example").WillReturnRows(
		sqlmock.NewRows([]string{"location"}).AddRow("http://example.com"),
	)
	mock.ExpectQuery("SELECT document FROM json_example").WillReturnRows(
		sqlmock.NewRows([]string{"document"}).AddRow("[1, 2, 3]"),
	)
	mock.ExpectQuery("SELECT id, name FROM scanner_example").WillReturnRows(
		sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "example").AddRow(2, "scanner"),
	)
	mock.ExpectQuery("SELECT id, name FROM into_struct_example").WillReturnRows(
		sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "example").AddRow(2, "scanner"),
	)
	mock.ExpectQuery(`SELECT name FROM builder_example WHERE id = \? AND found = \?`).WithArgs(1, true).WillReturnRows(
		sqlmock.NewRows([]string{"name"}).AddRow("example"),
	)
	mock.MatchExpectationsInOrder(false)
}

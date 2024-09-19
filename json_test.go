package sqlh_test

import (
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/shoenig/test/must"
	"github.com/simon-engledew/sqlh"
	"testing"
)

func ExampleJson() {
	var document any
	_ = db.QueryRow("SELECT document FROM json_example").Scan(sqlh.Json(&document))
	fmt.Println(document)
	// Output: [1 2 3]
}

func TestJson(t *testing.T) {
	db, mock, err := sqlmock.New()
	must.NoError(t, err)

	mock.ExpectQuery("SELECT id, document FROM test").WillReturnRows(
		sqlmock.NewRows([]string{"id", "document"}).AddRow(1, "[1, 2, 3]"),
	)

	var id int
	var document []int

	must.NoError(t, db.QueryRow("SELECT id, document FROM test").Scan(&id, sqlh.Json(&document)))

	must.EqOp(t, id, 1)
	must.SliceEqOp(t, document, []int{1, 2, 3})
}

func TestJsonRows(t *testing.T) {
	db, mock, err := sqlmock.New()
	must.NoError(t, err)
	t.Cleanup(func() {
		mock.ExpectClose()
		must.NoError(t, db.Close())
	})

	mock.ExpectQuery("SELECT document FROM test").WillReturnRows(
		sqlmock.NewRows([]string{"document"}).AddRow("[1, 2, 3]").AddRow("[2, 3, 4]"),
	)

	rows, err := sqlh.SQL(`SELECT document FROM test`).Query(db)
	must.NoError(t, err)

	documents, err := sqlh.ScanV(rows, func(v *[]int32, rows sqlh.Row) error {
		return rows.Scan(sqlh.Json(v))
	})
	must.NoError(t, err)
	must.Eq(t, documents, [][]int32{{1, 2, 3}, {2, 3, 4}})
}

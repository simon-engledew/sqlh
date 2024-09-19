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
	t.Cleanup(func() {
		_ = db.Close()
	})

	mock.ExpectQuery("SELECT id, document FROM test").WillReturnRows(
		sqlmock.NewRows([]string{"id", "document"}).AddRow(1, "[1, 2, 3]"),
	)

	var id int
	var document []int

	must.NoError(t, db.QueryRow("SELECT id, document FROM test").Scan(&id, sqlh.Json(&document)))

	must.EqOp(t, id, 1)
	must.SliceEqOp(t, document, []int{1, 2, 3})
}

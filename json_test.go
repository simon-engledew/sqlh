package sqlh_test

import (
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/simon-engledew/sqlh"
	"github.com/stretchr/testify/require"
	"testing"
)

func ExampleJson() {
	var document any
	_ = db().QueryRow("SELECT document FROM test").Scan(sqlh.Json(&document))
	fmt.Println(document)
	// Output: [1 2 3]
}

func TestJson(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = db.Close()
	})

	mock.ExpectQuery("SELECT id, document FROM test").WillReturnRows(
		sqlmock.NewRows([]string{"id", "document"}).AddRow(1, "[1, 2, 3]"),
	)

	var id int
	var document []int

	require.NoError(t, db.QueryRow("SELECT id, document FROM test").Scan(&id, sqlh.Json(&document)))

	require.Equal(t, id, 1)
	require.Equal(t, document, []int{1, 2, 3})
}

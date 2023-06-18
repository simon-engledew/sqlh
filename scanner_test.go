package sqlh_test

import (
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/simon-engledew/sqlh"
	"github.com/stretchr/testify/require"
	"testing"
)

func ExampleScanner() {
	scanner := sqlh.Scanner(func(item *struct {
		id   int
		name string
	}, scan func(...any) error) error {
		return scan(&item.id, &item.name)
	})

	items, _ := scanner(db.Query("SELECT id, name FROM scanner_example"))
	for _, item := range items {
		fmt.Println(item.id, item.name)
	}
	// Output: 1 example
	// 2 scanner
}

type testRow struct {
	id   int
	name string
}

func TestScanner(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = db.Close()
	})

	mock.ExpectQuery("SELECT id, name FROM test").WillReturnRows(
		sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "a").AddRow(2, "b"),
	)

	scan := sqlh.Scanner(func(row *testRow, scan func(...any) error) error {
		return scan(&row.id, &row.name)
	})

	expected := []*testRow{
		{1, "a"},
		{2, "b"},
	}

	rows, err := scan(db.Query("SELECT id, name FROM test"))
	require.NoError(t, err)
	require.Equal(t, expected, rows)
}

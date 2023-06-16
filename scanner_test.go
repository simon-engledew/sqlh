package sqlh_test

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/simon-engledew/sqlh"
	"github.com/stretchr/testify/require"
	"testing"
)

type row struct {
	id   int
	name string
}

func TestScanner(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	mock.ExpectQuery("SELECT id, name FROM test").WillReturnRows(
		sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "a").AddRow(2, "b"),
	)

	scan := sqlh.Scanner(func(row *row, scan func(...any) error) error {
		return scan(&row.id, &row.name)
	})

	expected := []*row{
		{1, "a"},
		{2, "b"},
	}

	rows, err := scan(db.Query("SELECT id, name FROM test"))
	require.NoError(t, err)
	require.Equal(t, expected, rows)
}

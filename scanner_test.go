package sqlh_test

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/simon-engledew/sqlh"
	"github.com/stretchr/testify/require"
	"net/url"
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

func TestJson(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	mock.ExpectQuery("SELECT id, name FROM test").WillReturnRows(
		sqlmock.NewRows([]string{"id", "items"}).AddRow(1, "[1, 2, 3]"),
	)

	var id int
	var items []int

	require.NoError(t, db.QueryRow("SELECT id, name FROM test").Scan(&id, sqlh.Json(&items)))

	require.Equal(t, id, 1)
	require.Equal(t, items, []int{1, 2, 3})
}

func TestBinary(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	mock.ExpectQuery("SELECT id, location FROM test").WillReturnRows(
		sqlmock.NewRows([]string{"id", "location"}).AddRow(1, "http://example.com"),
	)

	var id int
	var location url.URL

	require.NoError(t, db.QueryRow("SELECT id, location FROM test").Scan(&id, sqlh.Binary(&location)))

	require.Equal(t, id, 1)
	require.Equal(t, location, url.URL{
		Scheme: "http",
		Host:   "example.com",
	})
}

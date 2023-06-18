package sqlh_test

import (
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/simon-engledew/sqlh"
	"github.com/stretchr/testify/require"
	"net/url"
	"testing"
)

func ExampleBinary() {
	var location url.URL
	_ = db.QueryRow("SELECT location FROM binary_example").Scan(sqlh.Binary(&location))
	fmt.Println(location.String())
	// Output: http://example.com
}

func TestBinary(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = db.Close()
	})

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

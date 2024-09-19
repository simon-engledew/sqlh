package sqlh_test

import (
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/shoenig/test/must"
	"github.com/simon-engledew/sqlh"
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
	must.NoError(t, err)

	mock.ExpectQuery("SELECT id, location_string, location_bytes FROM test WHERE location = ?").WithArgs([]byte("http://q.com")).WillReturnRows(
		sqlmock.NewRows([]string{"id", "location_string", "location_bytes"}).AddRow(1, "http://example.com", []byte("http://example.com")),
	)

	var id int
	var locationA, locationB url.URL

	must.NoError(t, db.QueryRow("SELECT id, location_string, location_bytes FROM test WHERE location = ?", sqlh.Binary(&url.URL{Scheme: "http", Host: "q.com"})).Scan(&id, sqlh.Binary(&locationA), sqlh.Binary(&locationB)))

	must.EqOp(t, id, 1)
	must.EqOp(t, locationA, url.URL{
		Scheme: "http",
		Host:   "example.com",
	})
	must.EqOp(t, locationA, locationB)
}

func TestNotBinary(t *testing.T) {
	db, mock, err := sqlmock.New()
	must.NoError(t, err)

	mock.ExpectQuery("SELECT location FROM test").WillReturnRows(
		sqlmock.NewRows([]string{"location"}).AddRow(1),
	)

	var location url.URL
	must.Error(t, db.QueryRow("SELECT location FROM test").Scan(sqlh.Binary(&location)))
}

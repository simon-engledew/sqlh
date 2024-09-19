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

	mock.ExpectQuery("SELECT id, location FROM test").WillReturnRows(
		sqlmock.NewRows([]string{"id", "location"}).AddRow(1, "http://example.com"),
	)

	var id int
	var location url.URL

	must.NoError(t, db.QueryRow("SELECT id, location FROM test").Scan(&id, sqlh.Binary(&location)))

	must.EqOp(t, id, 1)
	must.EqOp(t, location, url.URL{
		Scheme: "http",
		Host:   "example.com",
	})
}

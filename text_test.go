package sqlh_test

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/shoenig/test/must"
	"github.com/simon-engledew/sqlh"
	"net"
	"testing"
)

func TestText(t *testing.T) {
	db, mock, err := sqlmock.New()
	must.NoError(t, err)

	mock.ExpectQuery("SELECT addr FROM test").WillReturnRows(
		sqlmock.NewRows([]string{"addr"}).AddRow("127.0.0.1"),
	)

	var addr net.IP
	must.NoError(t, db.QueryRow("SELECT addr FROM test").Scan(sqlh.Text(&addr)))

	must.Eq(t, net.IPv4(127, 0, 0, 1), addr)
}

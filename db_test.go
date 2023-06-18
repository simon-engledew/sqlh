package sqlh_test

import (
	"database/sql"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"runtime"
)

func db() *sql.DB {
	db, mock, err := sqlmock.New()
	if err != nil {
		panic(err)
	}

	pc, _, _, ok := runtime.Caller(1)
	if !ok {
		panic("could not get caller")
	}
	fn := runtime.FuncForPC(pc)
	switch fn.Name() {
	case "github.com/simon-engledew/sqlh_test.ExampleBinary":
		mock.ExpectQuery("SELECT location FROM test").WillReturnRows(
			sqlmock.NewRows([]string{"location"}).AddRow("http://example.com"),
		)
	case "github.com/simon-engledew/sqlh_test.ExampleJson":
		mock.ExpectQuery("SELECT document FROM test").WillReturnRows(
			sqlmock.NewRows([]string{"document"}).AddRow("[1, 2, 3]"),
		)
	case "github.com/simon-engledew/sqlh_test.ExampleScanner":
		mock.ExpectQuery("SELECT id, name FROM test").WillReturnRows(
			sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "example").AddRow(2, "scanner"),
		)
	default:
		panic(fmt.Errorf("unknown function %q", fn.Name()))
	}

	return db
}

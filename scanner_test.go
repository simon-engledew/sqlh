package sqlh_test

import (
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/shoenig/test/must"
	"github.com/simon-engledew/sqlh"
	"testing"
)

func ExampleScanner() {
	rows, _ := db.Query("SELECT id, name FROM scanner_example")
	items, _ := sqlh.Scan(rows, func(item *testRow, rows sqlh.Row) error {
		return rows.Scan(&item.id, &item.name)
	})
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
	must.NoError(t, err)
	t.Cleanup(func() {
		_ = db.Close()
	})

	mock.ExpectQuery("SELECT id, name FROM test").WillReturnRows(
		sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "a").AddRow(2, "b"),
	)

	rows, err := db.Query("SELECT id, name FROM test")
	must.NoError(t, err)

	items, err := sqlh.Scan(rows, func(row *testRow, rows sqlh.Row) error {
		return rows.Scan(&row.id, &row.name)
	})
	must.NoError(t, err)

	expected := []*testRow{
		{1, "a"},
		{2, "b"},
	}

	must.Eq(t, expected, items)
}

func TestPluck(t *testing.T) {
	db, mock, err := sqlmock.New()
	must.NoError(t, err)
	t.Cleanup(func() {
		_ = db.Close()
	})

	mock.ExpectQuery("SELECT id FROM pluck_example").WillReturnRows(
		sqlmock.NewRows([]string{"id"}).AddRow(1).AddRow(2),
	)

	expected := []uint64{1, 2}

	ids, err := sqlh.Pluck[uint64](db.Query("SELECT id FROM pluck_example"))

	must.NoError(t, err)
	must.SliceEqOp(t, expected, ids)
}

func TestScannerAnonymous(t *testing.T) {
	db, mock, err := sqlmock.New()
	must.NoError(t, err)
	t.Cleanup(func() {
		_ = db.Close()
	})

	mock.ExpectQuery("SELECT id, name FROM test").WillReturnRows(
		sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "a").AddRow(2, "b"),
	)

	rows, err := db.Query("SELECT id, name FROM test")
	must.NoError(t, err)

	items, err := sqlh.Scan(rows, func(row *struct {
		id   int
		name string
	}, rows sqlh.Row) error {
		return rows.Scan(&row.id, &row.name)
	})
	must.NoError(t, err)

	must.Len(t, 2, items)
	must.EqOp(t, 1, items[0].id)
	must.EqOp(t, 2, items[1].id)
}

package sqlh_test

import (
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/simon-engledew/sqlh"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func ExampleIntoStruct() {
	rows, _ := db.Query("SELECT id, name FROM scanner_example")
	items, _ := sqlh.Scan(rows, sqlh.IntoStruct[struct {
		ID   int
		Name string
	}](sqlh.FieldMatcher))
	for _, item := range items {
		fmt.Println(item.ID, item.Name)
	}
	// Output: 1 example
	// 2 scanner
}

type testStruct struct {
	ID        int
	FirstName string
	CreatedAt time.Time
}

func TestScanIntoWithGuess(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = db.Close()
	})

	now := time.Now()

	mock.ExpectQuery("SELECT id, first_name, created_at FROM test").WillReturnRows(
		sqlmock.NewRows([]string{"id", "first_name", "created_at"}).AddRow(1, "a", now).AddRow(2, "b", now),
	)

	rows, err := db.Query("SELECT id, first_name, created_at FROM test")
	require.NoError(t, err)

	items, err := sqlh.Scan(rows, sqlh.IntoStruct[testStruct](sqlh.FieldMatcher))
	require.NoError(t, err)

	require.Len(t, items, 2)
	require.Equal(t, 1, items[0].ID)
	require.Equal(t, "a", items[0].FirstName)
	require.Equal(t, 2, items[1].ID)
	require.Equal(t, "b", items[1].FirstName)
	require.Equal(t, now, items[1].CreatedAt)
}

type taggedTestStruct struct {
	ID   int    `json:"id"`
	Name string `json:"name,omitempty"`
}

func TestScanIntoWithTags(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = db.Close()
	})

	mock.ExpectQuery("SELECT id, name FROM test").WillReturnRows(
		sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "a").AddRow(2, "b"),
	)

	rows, err := db.Query("SELECT id, name FROM test")
	require.NoError(t, err)

	items, err := sqlh.Scan(rows, sqlh.IntoStruct[taggedTestStruct](sqlh.TagMatcher("json")))
	require.NoError(t, err)

	require.Len(t, items, 2)
	require.Equal(t, 1, items[0].ID)
	require.Equal(t, "a", items[0].Name)
	require.Equal(t, 2, items[1].ID)
	require.Equal(t, "b", items[1].Name)
}

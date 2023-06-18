package sqlh_test

import (
	"fmt"
	"github.com/simon-engledew/sqlh"
	"github.com/stretchr/testify/require"
	"testing"
)

func ExampleIn() {
	ids := []int{1, 2, 3}
	query := sqlh.SQL(`SELECT name FROM in_example WHERE id IN (?)`, sqlh.In(ids))
	fmt.Println(query.Statement, query.Args)
	// Output: SELECT name FROM in_example WHERE id IN (?, ?, ?) [1 2 3]
}

func ExampleSQL() {
	clause := sqlh.SQL("found = ?", true)
	expr := sqlh.SQL(`SELECT name FROM builder_example WHERE id = ? AND ?`, 1, clause)
	var name string
	_ = db.QueryRow(expr.Statement, expr.Args...).Scan(&name)
	fmt.Println(name)
	// Output: example
}

func TestSQL(t *testing.T) {
	a := sqlh.SQL(`SELECT 1 FROM a WHERE id = ?`, 1)

	b := sqlh.SQL(`SELECT 1 FROM b WHERE id = ?`, 2)

	c := sqlh.SQL(`SELECT * FROM (?) AS a, (?) AS b LIMIT ?, ?`, a, b, 1, 10)

	require.Equal(t, []any{1, 2, 1, 10}, c.Args)
	require.Equal(t, `SELECT * FROM (SELECT 1 FROM a WHERE id = ?) AS a, (SELECT 1 FROM b WHERE id = ?) AS b LIMIT ?, ?`, c.Statement)

	d := sqlh.SQL(`SELECT * FROM test WHERE id IN (?, ?, ?, ?)`, 1, 2, 3)

	require.Equal(t, []any{1, 2, 3}, d.Args)
	require.Equal(t, `SELECT * FROM test WHERE id IN (?, ?, ?, ?)`, d.Statement)

	e := sqlh.SQL(`SELECT * FROM test WHERE id IN (?)`, sqlh.In([]int{1, 2, 3}))

	require.Equal(t, []any{1, 2, 3}, e.Args)
	require.Equal(t, `SELECT * FROM test WHERE id IN (?, ?, ?)`, e.Statement)
}

package sqlh_test

import (
	"fmt"
	"github.com/simon-engledew/sqlh"
	"github.com/stretchr/testify/require"
	"testing"
)

func ExampleIn() {
	query := sqlh.SQL(`SELECT name FROM in_example WHERE id IN (?)`, sqlh.In([]int{1, 2, 3}))
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

func TestIn(t *testing.T) {
	for _, v := range []struct {
		Args     []int
		Expected string
	}{
		{[]int{1, 2, 3, 4, 5}, "?, ?, ?, ?, ?"},
		{[]int{1, 2, 3, 4}, "?, ?, ?, ?"},
		{[]int{1, 2, 3}, "?, ?, ?"},
		{[]int{1, 2}, "?, ?"},
		{[]int{1}, "?"},
		{[]int{}, ""},
	} {
		q := sqlh.In(v.Args)
		require.Equal(t, v.Expected, q.Statement)
		require.Len(t, q.Args, len(v.Args))
		for n := range v.Args {
			require.Equal(t, v.Args[n], q.Args[n])
		}
	}
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

	f := sqlh.SQL("(SELECT 1)")
	g := sqlh.SQL("(SELECT 2)")

	h := sqlh.SQL("SELECT * FROM test WHERE id IN (?)", sqlh.In([]any{f, g}))

	require.Equal(t, "SELECT * FROM test WHERE id IN ((SELECT 1), (SELECT 2))", h.Statement)
	require.Len(t, h.Args, 0)

	i := sqlh.SQL(`SELECT 1 FROM a WHERE id = ?`)
	require.Equal(t, `SELECT 1 FROM a WHERE id = ?`, i.Statement)

	j := sqlh.SQL(`INSERT INTO a (id, name) VALUES ?`, sqlh.Values(
		[]any{1, "hello"},
		[]any{2, "test"},
	))
	require.Equal(t, `INSERT INTO a (id, name) VALUES (?, ?), (?, ?)`, j.Statement)
	require.Equal(t, []any{1, "hello", 2, "test"}, j.Args)
}

func TestDebugSQL(t *testing.T) {
	a := sqlh.DebugSQL(`SELECT 1 FROM a WHERE id = ?`, 1)

	b := sqlh.DebugSQL(`SELECT 1 FROM b WHERE id = ?`, 2)

	sqlh.DebugSQL(`SELECT * FROM (?) AS a, (?) AS b LIMIT ?, ?`, a, b, 1, 10)
}

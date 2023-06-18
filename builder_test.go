package sqlh_test

import (
	"fmt"
	"github.com/simon-engledew/sqlh"
	"github.com/stretchr/testify/require"
	"testing"
)

func ExampleIn() {
	ids := []int{1, 2, 3}
	fmt.Println(sqlh.SQL(`SELECT name FROM a WHERE id IN (?)`, sqlh.In(ids)).Statement)
	// Output: SELECT name FROM a WHERE id IN (?, ?, ?)
}

func ExampleBuilder() {
	clause := sqlh.SQL("found = ?", true)
	query := sqlh.SQL(`SELECT name FROM a WHERE id = ? AND ?`, 1, clause)
	fmt.Println(query.Statement)
	fmt.Println(query.Args)
	// Output: SELECT name FROM a WHERE id = ? AND found = ?
	// [1 true]
}

func TestBuilder(t *testing.T) {
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

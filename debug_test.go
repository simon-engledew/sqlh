package sqlh_test

import (
	"fmt"
	"github.com/shoenig/test/must"
	"github.com/simon-engledew/sqlh"
	"testing"
)

func ExampleDebugSQL() {
	subquery := sqlh.DebugSQL(`SELECT id FROM users WHERE suspended_at IS NULL AND parent_id = ?`, 10)
	query := sqlh.DebugSQL(`SELECT event FROM events WHERE user_id IN (?)`, subquery)
	fmt.Println(query.Statement)
	// Output: /* debug_test.go:12 */ SELECT event FROM events WHERE user_id IN (
	//	/* debug_test.go:11 */ SELECT id FROM users WHERE suspended_at IS NULL AND parent_id = ?
	//)
}

func TestDebugSQL(t *testing.T) {
	must.StrContains(t, sqlh.DebugSQL(`SELECT event FROM events WHERE ?`, sqlh.SQL(`id = ?`, 1)).Statement, "debug_test.go")
}

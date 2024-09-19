package sqlh_test

import (
	"fmt"
	"github.com/simon-engledew/sqlh"
)

func ExampleDebugSQL() {
	subquery := sqlh.DebugSQL(`SELECT id FROM users WHERE suspended_at IS NULL AND parent_id = ?`, 10)
	query := sqlh.DebugSQL(`SELECT event FROM events WHERE user_id IN (?)`, subquery)
	fmt.Println(query.Statement)
	// Output: /* debug_test.go:10 */ SELECT event FROM events WHERE user_id IN (
	//	/* debug_test.go:9 */ SELECT id FROM users WHERE suspended_at IS NULL AND parent_id = ?
	//)
}

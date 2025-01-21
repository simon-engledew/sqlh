package sqlh_test

import (
	"github.com/shoenig/test/must"
	"github.com/simon-engledew/sqlh"
	"testing"
)

func TestExprString(t *testing.T) {
	stmt := sqlh.SQL(`SELECT test FROM data WHERE ?`, sqlh.SQL(`id = ?`, 1))
	must.EqOp(t, "SELECT test FROM data WHERE id = ?", stmt.String())
}

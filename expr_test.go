package sqlh_test

import (
	"fmt"
	"github.com/shoenig/test/must"
	"github.com/simon-engledew/sqlh"
	"testing"
)

func TestExprString(t *testing.T) {
	stmt := sqlh.SQL(`SELECT test FROM data WHERE ?`, sqlh.SQL(`id = ?`, 1))
	must.EqOp(t, "SELECT test FROM data WHERE id = ?", fmt.Sprintf("%s", stmt))
}

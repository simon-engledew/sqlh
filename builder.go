package sqlh

import (
	"strings"
)

type Expr struct {
	Statement string
	Args      []any
}

func (e *Expr) String() string {
	return e.Statement
}

// In takes parameters and returns an Expr that can be used in an SQL IN clause.
func In(items ...any) *Expr {
	var stmt string
	switch len(items) {
	case 0:
		stmt = ""
	case 1:
		stmt = "?"
	default:
		stmt = strings.Repeat(", ?", len(items))[2:]
	}
	return SQL(stmt, items...)
}

// SQL takes an SQL fragment and returns an Expr that flattens any nested Expr structs and their
// arguments.
func SQL(stmt string, args ...any) *Expr {
	expr := &Expr{}
	expr.Args = make([]any, 0, len(args))

	sections := strings.Split(stmt, "?")

	out := make([]string, 0, len(sections))

	for idx, section := range sections {
		out = append(out, section)
		if idx < len(args) {
			arg := args[idx]
			if subquery, ok := arg.(*Expr); ok {
				out = append(out, subquery.Statement)
				expr.Args = append(expr.Args, subquery.Args...)
				continue
			}
			expr.Args = append(expr.Args, arg)
		}
		out = append(out, "?")
	}

	expr.Statement = strings.Join(out[:len(out)-1], "")
	return expr
}

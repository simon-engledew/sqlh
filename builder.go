package sqlh

import (
	"context"
	"strings"
)

type Expr struct {
	Statement string
	Args      []any
}

func (e *Expr) Query() (string, []any) {
	return e.Statement, e.Args
}

func (e *Expr) QueryContext(ctx context.Context) (context.Context, string, []any) {
	return ctx, e.Statement, e.Args
}

func (e *Expr) String() string {
	return e.Statement
}

func In[T any](items []T) *Expr {
	args := make([]any, len(items))
	for n, item := range items {
		args[n] = item
	}
	var stmt string
	switch len(args) {
	case 0:
		stmt = ""
	case 1:
		stmt = "?"
	default:
		stmt = strings.Repeat(", ?", len(args))[2:]
	}
	return &Expr{
		Statement: stmt,
		Args:      args,
	}
}

func SQL(doc string, args ...any) *Expr {
	expr := &Expr{}
	expr.Args = make([]any, 0, len(args))

	sections := strings.Split(doc, "?")

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

package sqlh

import (
	"strings"
)

// In takes parameters and returns an Expr that can be used in an SQL IN clause.
func In[T any, S ~[]T](items S) Expr {
	switch size := len(items); size {
	case 0:
		return Expr{}
	case 1:
		return SQL("?", items[0])
	case 2:
		return SQL("?, ?", items[0], items[1])
	case 3:
		return SQL("?, ?, ?", items[0], items[1], items[2])
	default:
		var b strings.Builder
		b.Grow((size * 3) - 2)
		b.WriteString("?")

		args := make([]any, 1, len(items))
		args[0] = items[0]

		for _, item := range items[1:] {
			args = append(args, item)
			b.WriteString(", ?")
		}

		return SQL(b.String(), args...)
	}
}

// Values allows you to build a multi-row insert statement.
func Values(values ...[]any) Expr {
	items := make([]Expr, len(values))
	for n, value := range values {
		items[n] = SQL("(?)", In(value))
	}
	return In(items)
}

// SQL takes an SQL fragment and returns an Expr that flattens any nested queries and their
// arguments.
func SQL(stmt string, args ...any) Expr {
	if len(args) == 0 {
		return Expr{Statement: stmt}
	}

	var expr Expr

	stmtSize := len(stmt)
	argsSize := len(args)
	for _, arg := range args {
		if sub, ok := arg.(Expr); ok {
			stmtSize += len(sub.Statement)
			argsSize += len(sub.Args)
		}
	}

	expr.Args = make([]any, 0, argsSize)

	var b strings.Builder
	b.Grow(stmtSize)

	var end, start int
	for i := 0; i < len(args); i += 1 {
		idx := strings.IndexByte(stmt[end:], '?')
		if idx < 0 {
			break
		}

		start, end = end, end+idx+1

		arg := args[i]
		if sub, ok := arg.(Expr); ok {
			b.WriteString(stmt[start : end-1])
			b.WriteString(sub.Statement)

			expr.Args = append(expr.Args, sub.Args...)
		} else {
			b.WriteString(stmt[start:end])

			expr.Args = append(expr.Args, arg)
		}
	}

	b.WriteString(stmt[end:])

	expr.Statement = b.String()
	return expr
}

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
		params := make([]any, len(items))
		for i := range items {
			params[i] = items[i]
		}
		return SQL("?"+strings.Repeat(", ?", len(items)-1), params...)
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

	subqueries := 0
	stmtSize := len(stmt)
	argsSize := 0
	for _, arg := range args {
		if sub, ok := arg.(Expr); ok {
			subqueries += 1
			stmtSize += len(sub.Statement) - 1 // do not include space for the original ?
			argsSize += len(sub.Args)
		} else {
			argsSize += 1
		}
	}

	if subqueries == 0 {
		return Expr{Statement: stmt, Args: args}
	}

	var expr Expr

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

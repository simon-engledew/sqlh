package sqlh

import (
	"context"
	"database/sql"
	"fmt"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

type Expr struct {
	Statement string
	Args      []any
}

func (e Expr) String() string {
	return e.Statement
}

// Query calls db.Query, passing in the SQL statement and its arguments.
// See https://pkg.go.dev/database/sql#DB.Query
func (e Expr) Query(db interface {
	Query(query string, args ...any) (*sql.Rows, error)
}) (*sql.Rows, error) {
	return db.Query(e.Statement, e.Args...)
}

// QueryContext calls db.QueryContext, passing in the SQL statement and its arguments.
// See https://pkg.go.dev/database/sql#DB.QueryContext
func (e Expr) QueryContext(ctx context.Context, db interface {
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
}) (*sql.Rows, error) {
	return db.QueryContext(ctx, e.Statement, e.Args...)
}

// QueryRow calls db.QueryRow, passing in the SQL statement and its arguments.
// See https://pkg.go.dev/database/sql#DB.QueryRow
func (e Expr) QueryRow(db interface {
	QueryRow(query string, args ...any) *sql.Row
}) *sql.Row {
	return db.QueryRow(e.Statement, e.Args...)
}

// QueryRowContext calls db.QueryRowContext, passing in the SQL statement and its arguments.
// See https://pkg.go.dev/database/sql#DB.QueryRowContext
func (e Expr) QueryRowContext(ctx context.Context, db interface {
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}) *sql.Row {
	return db.QueryRowContext(ctx, e.Statement, e.Args...)
}

// Exec calls db.Exec, passing in the SQL statement and its arguments.
// See https://pkg.go.dev/database/sql#DB.Exec
func (e Expr) Exec(db interface {
	Exec(query string, args ...any) (sql.Result, error)
}) (sql.Result, error) {
	return db.Exec(e.Statement, e.Args...)
}

// ExecContext calls db.ExecContext, passing in the SQL statement and its arguments.
// See https://pkg.go.dev/database/sql#DB.ExecContext
func (e Expr) ExecContext(ctx context.Context, db interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}) (sql.Result, error) {
	return db.ExecContext(ctx, e.Statement, e.Args...)
}

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

func Values(values ...[]any) Expr {
	items := make([]Expr, len(values))
	for n, value := range values {
		items[n] = SQL("(?)", In(value))
	}
	return In(items)
}

func indent(v string) string {
	if !strings.Contains(v, "\n") {
		return v
	}
	return "\n\t" + strings.Join(strings.Split(strings.TrimSpace(v), "\n"), "\n\t") + "\n"
}

func ignoreError[T []byte | string](s T, _ error) string {
	return string(s)
}

var moduleRoot = sync.OnceValue(func() string {
	return filepath.Dir(ignoreError(exec.Command("go", "env", "GOMOD").Output()))
})

// DebugSQL annotates the query with the caller and indents it if it contains a newline.
func DebugSQL(stmt string, args ...any) Expr {
	_, file, line, _ := runtime.Caller(1)

	for n, arg := range args {
		if subquery, ok := arg.(Expr); ok {
			args[n] = Expr{Statement: indent(subquery.Statement), Args: subquery.Args}
		}
	}

	if path := ignoreError(filepath.Rel(moduleRoot(), file)); path != "" {
		return SQL(fmt.Sprintf("\n/* %s:%d */ %s", path, line, stmt), args...)
	}

	return SQL(stmt, args...)
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

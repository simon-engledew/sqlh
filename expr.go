package sqlh

import (
	"context"
	"database/sql"
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

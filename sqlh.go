// Package sqlh provides lightweight sql helpers.
package sqlh

import "database/sql"

// Rows allows wrappers for sql.Rows to be passed to the scanning functions.
type Rows interface {
	Close() error
	Next() bool
	Err() error
	Scan(...any) error
	ColumnTypes() ([]*sql.ColumnType, error)
}

var _ Rows = &sql.Rows{}

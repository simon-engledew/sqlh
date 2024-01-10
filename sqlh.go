// Package sqlh provides lightweight sql helpers.
package sqlh

import "database/sql"

// Row allows wrappers for sql.Rows to be passed to the scanning functions.
type Row interface {
	Scan(...any) error
	ColumnTypes() ([]*sql.ColumnType, error)
}

var _ Row = &sql.Rows{}

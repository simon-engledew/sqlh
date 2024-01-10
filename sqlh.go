// Package sqlh provides lightweight sql helpers.
package sqlh

import "database/sql"

type Row interface {
	Scan(...any) error
	ColumnTypes() ([]*sql.ColumnType, error)
}

var _ Row = &sql.Rows{}

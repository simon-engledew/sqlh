// Package sqlh provides lightweight sql helpers.
package sqlh

type Rows interface {
	Close() error
	Next() bool
	Err() error
	Scan(...any) error
}

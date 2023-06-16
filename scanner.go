package scanner

import (
	"database/sql"
	"database/sql/driver"
	"encoding"
	"encoding/json"
	"fmt"
)

type binaryType struct {
	value interface {
		encoding.BinaryMarshaler
		encoding.BinaryUnmarshaler
	}
}

func Binary(v interface {
	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler
}) interface {
	sql.Scanner
	driver.Valuer
} {
	return binaryType{value: v}
}

func (b binaryType) Scan(val interface{}) error {
	bytes, ok := val.([]byte)
	if !ok {
		return fmt.Errorf("expected bytes, got %T", val)
	}
	return b.value.UnmarshalBinary(bytes)
}

func (b binaryType) Value() (driver.Value, error) {
	return b.value.MarshalBinary()
}

type jsonType[T any] struct {
	value T
}

func Json[T any](v T) interface {
	sql.Scanner
	driver.Valuer
} {
	return jsonType[T]{value: v}
}

func (b jsonType[T]) Scan(val interface{}) error {
	bytes, ok := val.([]byte)
	if !ok {
		return fmt.Errorf("expected bytes, got %T", val)
	}
	return json.Unmarshal(bytes, b.value)
}

func (b jsonType[T]) Value() (driver.Value, error) {
	data, err := json.Marshal(b.value)
	return string(data), err
}

// Scanner takes a function that can scan a given query into P and returns a function
// that can be given (*sql.Rows, error) and will return a list of P.
// The benefit of this approach is that it does not need to use reflection to fetch struct names.
func Scanner[P *V, V any](scan func(P, func(...any) error) error) func(rows *sql.Rows, queryErr error) ([]P, error) {
	return func(rows *sql.Rows, queryErr error) (out []P, err error) {
		if queryErr != nil {
			return out, queryErr
		}

		defer func() {
			rowsErr := rows.Close()
			if rowsErr != nil {
				if err == nil {
					err = fmt.Errorf("failed to close rows %w", rowsErr)
				}
			}
		}()

		for rows.Next() {
			if err := rows.Err(); err != nil {
				return out, err
			}

			var v V

			if err := scan(&v, rows.Scan); err != nil {
				return out, fmt.Errorf("failed to scan rows: %w", err)
			}

			out = append(out, &v)
		}

		return out, rows.Err()
	}
}

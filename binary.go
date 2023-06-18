package sqlh

import (
	"database/sql"
	"database/sql/driver"
	"encoding"
	"fmt"
)

type binaryType struct {
	value interface {
		encoding.BinaryMarshaler
		encoding.BinaryUnmarshaler
	}
}

// Binary converts to or from a binary value.
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
	switch data := val.(type) {
	case []byte:
		return b.value.UnmarshalBinary(data)
	case string:
		return b.value.UnmarshalBinary([]byte(data))
	default:
		return fmt.Errorf("expected bytes, got %T", val)
	}
}

func (b binaryType) Value() (driver.Value, error) {
	return b.value.MarshalBinary()
}

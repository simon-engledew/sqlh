package sqlh

import (
	"database/sql"
	"database/sql/driver"
	"encoding"
	"fmt"
)

type textType struct {
	value interface {
		encoding.TextMarshaler
		encoding.TextUnmarshaler
	}
}

// Text converts to or from a binary value.
func Text(v interface {
	encoding.TextMarshaler
	encoding.TextUnmarshaler
}) interface {
	sql.Scanner
	driver.Valuer
} {
	return textType{value: v}
}

func (b textType) Scan(val interface{}) error {
	switch data := val.(type) {
	case []byte:
		return b.value.UnmarshalText(data)
	case string:
		return b.value.UnmarshalText([]byte(data))
	default:
		return fmt.Errorf("expected bytes, got %T", val)
	}
}

func (b textType) Value() (driver.Value, error) {
	return b.value.MarshalText()
}

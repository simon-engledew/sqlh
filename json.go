package sqlh

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type jsonType[T any] struct {
	value T
}

// Json converts to or from a json value.
func Json[T any](v T) interface {
	sql.Scanner
	driver.Valuer
} {
	return jsonType[T]{value: v}
}

func (b jsonType[T]) Scan(val interface{}) error {
	switch data := val.(type) {
	case []byte:
		return json.Unmarshal(data, b.value)
	case string:
		return json.Unmarshal([]byte(data), b.value)
	default:
		return fmt.Errorf("expected bytes, got %T", val)
	}
}

func (b jsonType[T]) Value() (driver.Value, error) {
	data, err := json.Marshal(b.value)
	return string(data), err
}

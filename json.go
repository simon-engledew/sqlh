package sqlh

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type jsonType struct {
	value any
}

// Json converts to or from a json value.
func Json(v any) interface {
	sql.Scanner
	driver.Valuer
} {
	return jsonType{value: v}
}

func (b jsonType) Scan(val interface{}) error {
	switch data := val.(type) {
	case []byte:
		return json.Unmarshal(data, b.value)
	case string:
		return json.Unmarshal([]byte(data), b.value)
	default:
		return fmt.Errorf("expected bytes, got %T", val)
	}
}

func (b jsonType) Value() (driver.Value, error) {
	data, err := json.Marshal(b.value)
	return string(data), err
}

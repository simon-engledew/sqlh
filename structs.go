package sqlh

import (
	"fmt"
	"reflect"
	"strings"
)

type FieldPredicate func(field reflect.StructField) bool

func (pred FieldPredicate) Find(t reflect.Type) (int, bool) {
	for i := 0; i < t.NumField(); i++ {
		if pred(t.Field(i)) {
			return i, true
		}
	}
	return 0, false
}

func IntoStruct[V any, P *V](matcher func(col string) FieldPredicate) func(P, Rows) error {
	cache := map[string]int{}
	return func(p P, rows Rows) error {
		types, err := rows.ColumnTypes()
		if err != nil {
			return err
		}

		v := reflect.Indirect(reflect.ValueOf(p))

		args := make([]any, len(types))

		for i, c := range types {
			name := c.Name()
			idx, ok := cache[name]
			if !ok {
				idx, ok = matcher(name).Find(v.Type())
				if !ok {
					return fmt.Errorf("field %q not found", name)
				}
				cache[name] = idx
			}

			args[i] = v.Field(idx).Addr().Interface()
		}

		return rows.Scan(args...)
	}
}

func FieldMatcher(col string) FieldPredicate {
	return func(field reflect.StructField) bool {
		i, j := 0, 0

		for ; i < len(col) && j < len(field.Name); i++ {
			sr := col[i]
			tr := field.Name[j]

			if sr == '_' {
				continue
			}

			j += 1

			if sr == tr {
				continue
			}

			if tr < sr {
				tr, sr = sr, tr
			}

			if 'A' <= sr && sr <= 'Z' && tr == sr+'a'-'A' {
				continue
			}

			return false
		}

		return i >= j
	}
}

func TagMatcher(key string) func(col string) FieldPredicate {
	return func(col string) FieldPredicate {
		return func(field reflect.StructField) bool {
			tag := field.Tag.Get(key)
			if i := strings.Index(tag, ","); i >= 0 {
				tag = tag[:i]
			}
			return tag == col
		}
	}
}

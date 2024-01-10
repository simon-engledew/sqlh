package sqlh

import (
	"fmt"
	"reflect"
	"strings"
)

type FieldPredicate = func(field reflect.StructField) bool

func IntoStruct[V any, P *V](matcher func(col string) FieldPredicate) func(P, Row) error {
	cache := map[string]int{}
	return func(p P, rows Row) error {
		columnTypes, err := rows.ColumnTypes()
		if err != nil {
			return err
		}

		valueOf := reflect.Indirect(reflect.ValueOf(p))
		typeOf := valueOf.Type()
		if typeOf.Kind() != reflect.Struct {
			return fmt.Errorf("IntoStruct of non-struct type %T", p)
		}

		args := make([]any, len(columnTypes))

		for i, columnType := range columnTypes {
			name := columnType.Name()
			idx, ok := cache[name]
			if !ok {
				pred := matcher(name)

				for ; idx < typeOf.NumField(); idx++ {
					if ok = pred(typeOf.Field(idx)); ok {
						break
					}
				}

				if !ok {
					return fmt.Errorf("field %q not found", name)
				}

				cache[name] = idx
			}

			args[i] = valueOf.Field(idx).Addr().Interface()
		}

		return rows.Scan(args...)
	}
}

func skipUnderscores(v string, i int) int {
	for ; i < len(v) && v[i] == '_'; i++ {
	}
	return i
}

func FieldMatcher(col string) FieldPredicate {
	return func(field reflect.StructField) bool {
		i, j := skipUnderscores(col, 0), 0

		for ; i < len(col) && j < len(field.Name); i, j = skipUnderscores(col, i+1), j+1 {
			sr := col[i]
			tr := field.Name[j]

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

		return i == len(col) && j == len(field.Name)
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

func FromStruct[T any](matcher func(col string) FieldPredicate, cols []string) (func(v T) []any, error) {
	cache := map[string]int{}

	{
		var zero T
		typeOf := reflect.TypeOf(zero)

		if typeOf.Kind() != reflect.Struct {
			return nil, fmt.Errorf("FromStruct of non-struct type %T", zero)
		}

		for _, name := range cols {
			idx, ok := cache[name]
			if !ok {
				pred := matcher(name)

				for ; idx < typeOf.NumField(); idx++ {
					if ok = pred(typeOf.Field(idx)); ok {
						break
					}
				}

				if !ok {
					return nil, fmt.Errorf("field %q not found", name)
				}

				cache[name] = idx
			}
		}
	}

	return func(v T) []any {
		valueOf := reflect.ValueOf(v)

		args := make([]any, len(cols))

		for i, col := range cols {
			args[i] = valueOf.Field(cache[col]).Interface()
		}

		return args
	}, nil
}

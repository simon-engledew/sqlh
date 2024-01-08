package sqlh

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

// Pluck will scan the results of a query that produces a single column.
func Pluck[V any](rows Rows, queryErr error) (out []V, err error) {
	if queryErr != nil {
		return out, queryErr
	}
	return ScanV(rows, func(v *V, rows Rows) error {
		return rows.Scan(v)
	})
}

// Iter calls fn for each successful call to rows.Next, closing the rows at the end.
func Iter(rows Rows, fn func() error) (err error) {
	defer func() {
		err = errors.Join(err, rows.Close())
	}()

	for rows.Next() {
		if err := rows.Err(); err != nil {
			return err
		}

		if err := fn(); err != nil {
			return fmt.Errorf("failed to scan rows: %w", err)
		}
	}

	return rows.Err()
}

// ScanV takes a function that can scan a given sql.Rows into []V.
func ScanV[P *V, V any](rows Rows, scan func(P, Rows) error) (out []V, err error) {
	err = Iter(rows, func() error {
		var v V
		err := scan(&v, rows)
		if err == nil {
			out = append(out, v)
		}
		return err
	})
	return
}

// Scan takes a function that can scan a given sql.Rows into []P.
func Scan[P *V, V any](rows Rows, scan func(P, Rows) error) (out []P, err error) {
	err = Iter(rows, func() error {
		var v V
		err := scan(&v, rows)
		if err == nil {
			out = append(out, &v)
		}
		return err
	})
	return
}

func (pred FieldPredicate) Find(t reflect.Type) (int, bool) {
	for i := 0; i < t.NumField(); i++ {
		if pred(t.Field(i)) {
			return i, true
		}
	}
	return 0, false
}

type StructMatcher func(col string) FieldPredicate

type FieldPredicate func(field reflect.StructField) bool

func IntoStruct[V any, P *V](matcher StructMatcher) func(P, Rows) error {
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

func IntoFields[V any, P *V]() func(P, Rows) error {
	return IntoStruct[V, P](FieldMatcher)
}

func FieldMatcher(col string) FieldPredicate {
	guess := strings.ReplaceAll(col, "_", "")
	return func(field reflect.StructField) bool {
		return strings.EqualFold(guess, field.Name)
	}
}

func IntoTag[V any, P *V](key string) func(P, Rows) error {
	return IntoStruct[V, P](TagMatcher(key))
}

func TagMatcher(key string) StructMatcher {
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

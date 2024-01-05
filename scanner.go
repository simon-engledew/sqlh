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

func guessField(t reflect.Type, name string) int {
	for i := 0; i < t.NumField(); i++ {
		if strings.EqualFold(name, t.Field(i).Name) {
			return i
		}
	}
	panic("field not found")
}

func ToStruct[V any, P *V]() func(P, Rows) error {
	fields := map[string]int{}

	{
		var zero V
		zeroType := reflect.TypeOf(zero)
		for i := 0; i < zeroType.NumField(); i++ {
			if field := zeroType.Field(i); field.IsExported() {
				if name := field.Tag.Get("sql"); name != "" {
					fields[name] = i
				}
			}
		}
	}

	return func(p P, rows Rows) error {
		types, err := rows.ColumnTypes()
		if err != nil {
			return err
		}

		v := reflect.Indirect(reflect.ValueOf(p))

		args := make([]any, len(types))

		for i, c := range types {
			idx, ok := fields[c.Name()]
			if !ok {
				idx = guessField(v.Type(), strings.ReplaceAll(c.Name(), "_", ""))
				fields[c.Name()] = idx
			}

			args[i] = v.Field(idx).Addr().Interface()
		}

		return rows.Scan(args...)
	}
}

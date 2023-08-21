package sqlh

import (
	"fmt"
)

func Pluck[V any](rows Rows, queryErr error) (out []V, err error) {
	if queryErr != nil {
		return out, queryErr
	}
	return ScanV(rows, func(v *V, scan func(...any) error) error {
		return scan(v)
	})
}

func Iter(rows Rows, fn func() error) (err error) {
	defer func() {
		rowsErr := rows.Close()
		if rowsErr != nil {
			if err == nil {
				err = fmt.Errorf("failed to close rows: %w", rowsErr)
			}
		}
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

func ScanV[P *V, V any](rows Rows, scan func(P, func(...any) error) error) (out []V, err error) {
	err = Iter(rows, func() error {
		var v V
		err := scan(&v, rows.Scan)
		if err == nil {
			out = append(out, v)
		}
		return err
	})
	return
}

// Scan takes a function that can scan a given sql.Rows into []P.
func Scan[P *V, V any](rows Rows, scan func(P, func(...any) error) error) (out []P, err error) {
	err = Iter(rows, func() error {
		var v V
		err := scan(&v, rows.Scan)
		if err == nil {
			out = append(out, &v)
		}
		return err
	})
	return
}

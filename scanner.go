package sqlh

import (
	"database/sql"
	"fmt"
)

func Pluck[V any](rows *sql.Rows, queryErr error) (out []V, err error) {
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

		if err := rows.Scan(&v); err != nil {
			return out, fmt.Errorf("failed to scan rows: %w", err)
		}

		out = append(out, v)
	}

	return out, rows.Err()
}

// Scan takes a function that can scan a given sql.Rows into []P.
func Scan[P *V, V any](rows *sql.Rows, scan func(P, func(...any) error) error) (out []P, err error) {
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
			return out, err
		}

		var v V

		if err := scan(&v, rows.Scan); err != nil {
			return out, fmt.Errorf("failed to scan rows: %w", err)
		}

		out = append(out, &v)
	}

	return
}

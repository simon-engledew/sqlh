package sqlh

import (
	"database/sql"
	"fmt"
)

// Scanner takes a function that can scan a given query into P and returns a function
// that can be given (*sql.Rows, error) and will return a list of P.
func Scanner[P *V, V any](scan func(P, func(...any) error) error) func(rows *sql.Rows, queryErr error) ([]P, error) {
	return func(rows *sql.Rows, queryErr error) (out []P, err error) {
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

			if err := scan(&v, rows.Scan); err != nil {
				return out, fmt.Errorf("failed to scan rows: %w", err)
			}

			out = append(out, &v)
		}

		return out, rows.Err()
	}
}

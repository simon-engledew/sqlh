[![Build Status](https://github.com/simon-engledew/sqlh/actions/workflows/tests.yml/badge.svg?branch=main)](https://github.com/simon-engledew/sqlh/actions)
[![GoDoc](https://godoc.org/github.com/simon-engledew/sqlh?status.svg)](https://godoc.org/github.com/simon-engledew/sqlh)
[![Go Report Card](https://goreportcard.com/badge/github.com/simon-engledew/sqlh)](https://goreportcard.com/report/github.com/simon-engledew/sqlh)

# sqlh

A lightweight Go database helper library.

Contains generic functions for scanning sql.Rows into structs and building SQL queries with IN clauses.

```go
package example

import (
	"time"
	"context"
	"fmt"
	"database/sql"
	"github.com/simon-engledew/sqlh"
)

type LatestEvent struct {
	name string
	updatedAt time.Time
}

func Example(ctx context.Context, db *sql.DB) error {
	// build a subquery with an argument
    grouped := sqlh.SQL(`SELECT events.name, MAX(events.updated_at) AS 'updated_at'
    FROM events
    WHERE events.owner_id = ?
    GROUP BY events.name`, 10)

	// combine the two queries
    query := sqlh.SQL(`SELECT grouped.name, grouped.updated_at
    FROM (?) AS grouped
    grouped.updated_at DESC
    LIMIT ?, ?`, grouped, 10, 100)

	// execute the query and get some rows back
    rows, err := db.QueryContext(ctx, query.Statement, query.Args...)
    if err != nil {
        return err
    }

	// use Scan to build a slice of structs from each row
    events, err := sqlh.Scan(rows, func(e *LatestEvent, scan func(dest ...any) error) error {
        return scan(&e.name, &e.updatedAt)
    })
	if err != nil {
		return err
    }

	for _, event := range events {
		fmt.Println(event.name, event.updatedAt)
    }

	return nil
}
```
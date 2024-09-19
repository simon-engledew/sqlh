[![Build Status](https://github.com/simon-engledew/sqlh/actions/workflows/tests.yml/badge.svg?branch=main)](https://github.com/simon-engledew/sqlh/actions)
[![GoDoc](https://godoc.org/github.com/simon-engledew/sqlh?status.svg)](https://godoc.org/github.com/simon-engledew/sqlh)
[![Go Report Card](https://goreportcard.com/badge/github.com/simon-engledew/sqlh)](https://goreportcard.com/report/github.com/simon-engledew/sqlh)

# sqlh

A lightweight Go database helper library.

Attempts to make SQL easier in Go without adding all the weight of an ORM.

## Problems this library is trying to solve

### I want to be able to use slice parameters for IN clauses!

```go
names, err := sqlh.Pluck[string](
  sqlh.SQL(`SELECT name FROM users WHERE id IN (?)`, sqlh.In(userIDs)).Query(db),
)
// SELECT name FROM users WHERE id IN (1, 2, 3)
```

### I want to be able to use SQL subqueries as parameters!

```go
subquery := sqlh.SQL(`SELECT id FROM users WHERE suspended_at IS NULL AND parent_id = ?`, 10)
query := sqlh.SQL(`SELECT event FROM events WHERE user_id IN (?)`, subquery)
// SELECT event FROM events WHERE user_id IN (SELECT id FROM users WHERE suspended_at IS NULL AND parent_id = ?), 10
```

### I want to be able to use slices as multi INSERT values!

```go
query := sqlh.SQL(`INSERT INTO a (id, name) VALUES ?`, sqlh.Values([]any{1, "hello"}, []any{2, "test"}))
// "INSERT INTO a (id, name) VALUES (?, ?), (?, ?)", [1, "hello", 2, "test"]
```

### I want to scan SQL values directly into protobuf structs!

```go
rows, err := sqlh.SQL(`SELECT name, address FROM users WHERE id IN (?)`, sqlh.In(userIDs)).Query(db)
users, err := sqlh.Scan[proto.UserResponse](rows, func (v *proto.UserResponse, row sqlh.Row) error {
	return row.Scan(&v.UserName, &v.Address)
})
```

### I want to scan JSON values directly into their go counterparts!

```go
var out map[string]string
err := sqlh.SQL(`SELECT json_data FROM store`).Scan(sqlh.Json(&out))
// {"hello": "test"}
```

### I want to scan JSON values as parameters!

```go
rows, err := sqlh.SQL(`SELECT json_data FROM store WHERE json_overlaps(json_data, ?)`, sqlh.Json(map[string]string {"hello": "test"})).Query(db)
// {"hello": "test"}
```

### Ok, but surely this isn't going to scale to a big codebase? I want to see where query fragments are coming from!

```go
subquery := sqlh.DebugSQL(`SELECT id FROM users WHERE suspended_at IS NULL AND parent_id = ?`, 10)
query := sqlh.DebugSQL(`SELECT event FROM events WHERE user_id IN (?)`, subquery)
// query is annotated with the place it was created:
// /* debug_test.go:10 */ SELECT event FROM events WHERE user_id IN (
//   /* debug_test.go:9 */ SELECT id FROM users WHERE suspended_at IS NULL AND parent_id = ?
// )
```

## Example

```go
var SQL = sqlh.SQL

if testing.Testing() {
    SQL = sqlh.DebugSQL
}

type LatestEvent struct {
	name string
	updatedAt time.Time
}

ownerIDs := []int{1, 5, 10}

// build a subquery with an argument
grouped := SQL(`SELECT events.name, MAX(events.updated_at) AS 'updated_at'
FROM events
WHERE events.owner_id IN (?)
GROUP BY events.name`, sqlh.In(ownerIDs))

// combine the two queries
query := SQL(`SELECT grouped.name, grouped.updated_at
FROM (?) AS grouped
grouped.updated_at DESC
LIMIT ?, ?`, grouped, 10, 100)

// execute the query and get some rows back
rows, err := query.QueryContext(ctx, db)
if err != nil { return err }

// use Scan to build a slice of structs from each row
events, err := sqlh.Scan(rows, func(e *LatestEvent, scan func(dest ...any) error) error {
	return scan(&e.name, &e.updatedAt)
})
if err != nil { return err }

for _, event := range events {
	fmt.Println(event.name, event.updatedAt)
}
```

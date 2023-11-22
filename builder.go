package sqlh

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

type Expr struct {
	Statement string
	Args      []any
}

func (e Expr) String() string {
	return e.Statement
}

// In takes parameters and returns an Expr that can be used in an SQL IN clause.
func In[T any, S ~[]T](items S) Expr {
	args := make([]any, 0, len(items))
	for _, item := range items {
		args = append(args, item)
	}
	var stmt string
	switch len(items) {
	case 0:
		stmt = ""
	case 1:
		stmt = "?"
	default:
		stmt = strings.Repeat(", ?", len(items))[2:]
	}
	return SQL(stmt, args...)
}

func indent(v string) string {
	if !strings.Contains(v, "\n") {
		return v
	}
	return "\n\t" + strings.Join(strings.Split(strings.TrimSpace(v), "\n"), "\n\t") + "\n"
}

var moduleRoot = sync.OnceValue(func() string {
	dir, _ := os.Getwd()
	if dir == "" {
		return dir
	}

	dir = filepath.Clean(dir)

	// Look for enclosing go.mod.
	for {
		if fi, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil && !fi.IsDir() {
			return dir
		}
		d := filepath.Dir(dir)
		if d == dir {
			break
		}
		dir = d
	}

	return ""
})

// DebugSQL annotates the query with the caller and indents it if it contains a newline.
func DebugSQL(stmt string, args ...any) Expr {
	_, file, line, _ := runtime.Caller(1)
	path, _ := filepath.Rel(moduleRoot(), file)

	for n, arg := range args {
		if subquery, ok := arg.(Expr); ok {
			args[n] = Expr{Statement: indent(subquery.Statement), Args: subquery.Args}
		}
	}

	return SQL(fmt.Sprintf("\n/* %s:%d */ %s", path, line, stmt), args...)
}

// SQL takes an SQL fragment and returns an Expr that flattens any nested queries and their
// arguments.
func SQL(stmt string, args ...any) Expr {
	var expr Expr
	expr.Args = make([]any, 0, len(args))

	sections := strings.Split(stmt, "?")

	out := make([]string, 0, len(sections))

	for idx, section := range sections {
		out = append(out, section)
		if idx < len(args) {
			arg := args[idx]
			if subquery, ok := arg.(Expr); ok {
				out = append(out, subquery.Statement)
				expr.Args = append(expr.Args, subquery.Args...)
				continue
			}
			expr.Args = append(expr.Args, arg)
		}
		out = append(out, "?")
	}

	expr.Statement = strings.Join(out[:len(out)-1], "")
	return expr
}

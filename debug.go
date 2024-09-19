package sqlh

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

func indent(v string) string {
	if !strings.Contains(v, "\n") {
		return v
	}
	return "\n\t" + strings.Join(strings.Split(strings.TrimSpace(v), "\n"), "\n\t") + "\n"
}

func ignoreError[T []byte | string](s T, _ error) string {
	return string(s)
}

var moduleRoot = sync.OnceValue(func() string {
	return filepath.Dir(ignoreError(exec.Command("go", "env", "GOMOD").Output()))
})

// DebugSQL annotates the query with the caller and indents it if it contains a newline.
func DebugSQL(stmt string, args ...any) Expr {
	_, file, line, _ := runtime.Caller(1)

	for n, arg := range args {
		if subquery, ok := arg.(Expr); ok {
			args[n] = Expr{Statement: indent(subquery.Statement), Args: subquery.Args}
		}
	}

	if path := ignoreError(filepath.Rel(moduleRoot(), file)); path != "" {
		return SQL(fmt.Sprintf("\n/* %s:%d */ %s", path, line, stmt), args...)
	}

	return SQL(stmt, args...)
}

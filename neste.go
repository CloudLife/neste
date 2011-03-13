/*
	Extended version of Go's template package for generating textual output 
	from nested templates.

	neste template engine makes it easier to generate output from multiple 
	nested template files by naming and managing all templates and also allowing
	generating output from them directly as strings.

	neste also includes many useful built-in formatters.
*/
package neste

import (
	"os"
)

// getMtime returns modified time of the given file.
func getMtime(path string) int64 {
	fi, _ := os.Lstat(path)
	return fi.Mtime_ns
}

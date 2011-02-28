// neste template engine: built-in formatters

package neste

import (
	"io"
	"template"
	"fmt"
	"bytes"
	"utf8"
	"unicode"
)

var builtins = template.FormatterMap{
	"e":          template.HTMLFormatter, // Just a shorthand for the "html" escaping formatter
	"addSlashes": AddSlashesFormatter,
	"capFirst":   CapFirstFormatter}

/*
Adds slashes before quotes. Useful for escaping strings in CSV, for example.

Example:

	{value|addslashes}

If value is "I'm using neste", the output will be "I\'m using neste".
*/
func AddSlashesFormatter(w io.Writer, formatter string, data ...interface{}) {
	b := getBytes(data...)

	for _, v := range b {
		if v == '"' {
			w.Write([]byte{'\\', '"'})
		} else {
			w.Write([]byte{v})
		}
	}
}

/*
Capitalizes the first character of the value.

Example:

	{value|capfirst}

If value is "neste", the output will be "Neste".
*/
func CapFirstFormatter(w io.Writer, formatter string, data ...interface{}) {
	b := getBytes(data...)

	if len(b) > 0 {
		rune, size := utf8.DecodeRune(b)
		rune = unicode.ToUpper(rune)
		capSize := utf8.RuneLen(rune)
		capb := make([]byte, capSize)
		utf8.EncodeRune(capb, rune)
		w.Write(capb)
		w.Write(b[size:])
	}
}

// Returns a byte slice of the (first) field value.
func getBytes(data ...interface{}) (b []byte) {
	ok := false

	if len(data) == 1 {
		b, ok = data[0].([]byte)
	}

	if !ok {
		var buf bytes.Buffer
		fmt.Fprint(&buf, data...)
		b = buf.Bytes()
	}
	return b
}


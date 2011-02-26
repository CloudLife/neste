// neste template engine: template

package neste

import (
	"template"
	"os"
	"bytes"
	"log"
	"io"
)

// Template is a type for holding a *template.Template.
type Template struct {
	cache *template.Template
}

// Execute applies a parsed template to the specified data object, generating output to wr.
// If any errors occur, err will be non-nil. 
func (t *Template) Execute(wr io.Writer, data interface{}) (err os.Error) {
	tt := t.cache
	err = tt.Execute(wr, data)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

// Render applies a parsed template to the specified data object and 
// returns the generated output as a string.
// If any errors occur, err will be non-nil. 
func (t *Template) Render(data interface{}) (s string, err os.Error) {
	buf := new(bytes.Buffer)
	err = t.Execute(buf, data)
	if err != nil {
		log.Println(err)
	}
	
	s = string(buf.Bytes())
	return s, nil
}


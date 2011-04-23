// neste template engine: template

package neste

import (
	"template"
	"os"
	"bytes"
	"io"
	"path"
)

type templateFileInfo struct {
	filename  string
	mtime     int64 // Modified time
	mustParse bool
}

// Template is a type for holding a *template.Template and other information.
type Template struct {
	m     *Manager
	cache *template.Template
	fi    *templateFileInfo // Used only for template files
}

// Execute applies a parsed template to the specified data object, 
// generating output to wr. If the template is a template file and the 
// template's template manager has reloading mode enabled, 
// then this method will attempt to reparse the template file if its modified 
// time has changed.
// If any errors occur, err will be non-nil.
func (t *Template) Execute(wr io.Writer, data interface{}) (err os.Error) {
	if t.fi != nil && t.m.reloading {
		err = t.Reload()
		if err != nil {
			return err
		}
	}

	tt := t.cache
	err = tt.Execute(wr, data)
	if err != nil {
		return err
	}

	return nil
}

// Reload rereads and reparses the template's associated template file
// if its modified time has changed since initial loading.
// Calling this method is unnecessary when reloading mode is enabled,
// unless the file's modified time is erroneous.
// If any errors occur, err will be non-nil.
func (t *Template) Reload() (err os.Error) {
	filename := t.fi.filename
	path := path.Join(t.m.baseDir, filename)
	oldMtime := t.fi.mtime
	curMtime := getMtime(path)

	if curMtime > oldMtime {
		// Template has changed.
		// Reparse the template file.
		t.cache, err = t.m.parsett(path, t.fi.mustParse)
		if err != nil {
			return err
		}
		
		// Update modified time
		t.fi.mtime = getMtime(path)
	}

	return nil
}

// Render applies a parsed template to the specified data object and 
// returns the generated output as a string.
// If any errors occur, output will be empty string "" and err will be non-nil. 
func (t *Template) Render(data interface{}) (s string, err os.Error) {
	buf := new(bytes.Buffer)
	err = t.Execute(buf, data)
	if err != nil {
		return "", err
	}

	s = string(buf.Bytes())
	return s, nil
}


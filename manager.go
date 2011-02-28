/*
	Extended version of Go's template package for generating textual output from nested templates.

	neste template engine makes it easier to generate output from multiple nested template files
	by naming and managing all templates and also allowing generating output from them directly as strings.

	neste also includes many useful built-in formatters.
*/
package neste

import (
	"template"
	"os"
	"log"
	"io/ioutil"
	"path"
)

// Manager is a type that represents a template manager.
type Manager struct {
	fmap     template.FormatterMap
	baseDir  string
	tStrings map[string]*Template // Templates for strings
	tFiles   map[string]*Template // Templates for files
	ldelim   string
	rdelim   string
}

// Returns a new template manager with base directory baseDir for template files.
func New(baseDir string, fmap template.FormatterMap) *Manager {
	// Add each built-in formatter unless there's a user given formatter with same name already.
	if fmap != nil {
		for k, v := range builtins {
			_, present := fmap[k]
			if !present {
				fmap[k] = v
			}
		}
	}

	return &Manager{
		baseDir:  baseDir,
		tStrings: make(map[string]*Template),
		tFiles:   make(map[string]*Template),
		fmap:     fmap,
		ldelim:   "{{",
		rdelim:   "}}"}
}

// Add adds a given template string s to the template manager with the identifier id.
// If any errors occur, returned error will be non-nil. 
func (tm *Manager) Add(s string, id string) (*Template, os.Error) {
	return tm.add(s, id, false)
}

// AddFile adds a given template file to the template manager.
// This method ignores ending newline character in the input file, unlike AddFilenl,
// to prevent accumulating extra newlines when nesting is performed. 
// If any errors occur, returned error will be non-nil. 
func (tm *Manager) AddFile(filename string) (*Template, os.Error) {
	return tm.addFile(filename, true, false)
}

// AddFile adds a given template file to the template manager.
// This method does not ignore ending newline character in the input file, unlike AddFile does
// to prevent accumulating extra newlines when nesting is performed. 
// If any errors occur, returned error will be non-nil. 
func (tm *Manager) AddFilenl(filename string) (*Template, os.Error) {
	return tm.addFile(filename, false, false)
}

// Removes all templates from the template manager.
// Useful for clearing out cached templates.
// Clear returns true if one or more templates were removed, otherwise false.
func (tm *Manager) Clear() bool {
	tlen := len(tm.tStrings) + len(tm.tFiles)
	tm.tStrings = make(map[string]*Template)
	tm.tFiles = make(map[string]*Template)
	return tlen > 0
}

// Returns a template with the given identifier or nil if it doesn't exist.
func (tm *Manager) Get(s string) *Template {
	return tm.tStrings[s]
}

// Returns a template with the given filename or nil if it doesn't exist.
func (tm *Manager) GetFile(filename string) *Template {
	path := path.Join(tm.baseDir, filename)
	return tm.tFiles[path]
}

// MustAdd is like Add, but panics, if template can't be parsed. 
func (tm *Manager) MustAdd(s string, id string) *Template {
	t, _ := tm.add(s, id, true)
	return t
}


// MustAddFile is like AddFile, but panics, if template can't be parsed. 
func (tm *Manager) MustAddFile(filename string) *Template {
	t, _ := tm.addFile(filename, true, true)
	return t
}


// MustAddFilenl is like AddFilenl, but panics, if template can't be parsed. 
func (tm *Manager) MustAddFilenl(filename string) *Template {
	t, _ := tm.addFile(filename, false, true)
	return t
}

// Removes a template with the given identifier from the template manager.
// Useful for clearing out cached templates.
// It's safe to remove a non-existing template.
// Remove returns true if a template was removed, otherwise false.
func (tm *Manager) Remove(s string) bool {
	_, present := tm.tStrings[s]
	tm.tStrings[s] = nil, false
	return present
}

// Removes a template with the given filename from the template manager.
// Useful for clearing out cached templates.
// It's safe to remove a non-existing template.
// Remove returns true if a template was removed, otherwise false.
func (tm *Manager) RemoveFile(filename string) bool {
	path := path.Join(tm.baseDir, filename)
	_, present := tm.tFiles[path]
	tm.tFiles[path] = nil, false
	return present
}

// SetDelims sets the left and right delimiters for operations in the template for template parsing.
func (tm *Manager) SetDelims(left, right string) {
	tm.ldelim = left
	tm.rdelim = right
}


// Unexported methods

// Add adds a given template string to the template manager.
// This method should not be called directly, but through Add or MustAdd.
// If any errors occur, err will be non-nil. 
func (tm *Manager) add(s string, id string, mustParse bool) (t *Template, err os.Error) {
	tt := template.New(tm.fmap)
	tt.SetDelims(tm.ldelim, tm.rdelim)

	// Parse the template.
	if mustParse {
		err := tt.Parse(s)
		if err != nil {
			log.Panic(err)
		}
	} else {
		err := tt.Parse(s)
		if err != nil {
			log.Println(err)
			return nil, err
		}
	}

	t = &Template{cache: tt}

	// Add template to the manager.
	tm.tStrings[id] = t
	return t, nil
}

// readFileNl is same as ioutil.ReadFile except it ignores ending newline character.
func readFileNl(path string) (in []byte, err os.Error) {
	in, err = ioutil.ReadFile(path)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	// Remove ending newline character only if one exists.
	if len(in) > 0 && in[len(in)-1] == '\n' {
		in = in[:len(in)-1]
	}
	return in, nil
}

// AddFile adds a given template file to the template manager.
// This method should not be called directly, but through AddFile, AddFilenl, MustAddFile or MustAddFilenl.
// ignoreEndingNl determines whether ending newline character is ignored in the input file.
// If any errors occur, err will be non-nil. 
func (tm *Manager) addFile(filename string, ignoreEndingNl bool, mustParse bool) (t *Template, err os.Error) {
	tt := template.New(tm.fmap)
	tt.SetDelims(tm.ldelim, tm.rdelim)

	// Parse template file.
	path := path.Join(tm.baseDir, filename)
	if ignoreEndingNl {
		if mustParse {
			var bstr []byte

			bstr, err = readFileNl(path)
			if err != nil {
				log.Panic(err)
			}

			err := tt.Parse(string(bstr))
			if err != nil {
				log.Panic(err)
			}
		} else {
			var bstr []byte

			bstr, err = readFileNl(path)
			if err != nil {
				log.Println(err)
				return nil, err
			}

			err := tt.Parse(string(bstr))
			if err != nil {
				log.Println(err)
				return nil, err
			}
		}
	} else {
		if mustParse {
			err := tt.ParseFile(path)
			if err != nil {
				log.Panic(err)
			}
		} else {
			err := tt.ParseFile(path)
			if err != nil {
				log.Println(err)
				return nil, err
			}
		}
	}

	t = &Template{cache: tt}

	// Add template to the manager.
	tm.tFiles[path] = t

	return t, nil
}


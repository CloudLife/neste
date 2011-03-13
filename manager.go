package neste

import (
	"template"
	"os"
	"path"
)

// Manager is a type that represents a template manager.
type Manager struct {
	fmap      template.FormatterMap
	baseDir   string
	tStrings  map[string]*Template // Templates for strings
	tFiles    map[string]*Template // Templates for files
	ldelim    string
	rdelim    string
	reloading bool
}

// Returns a new template manager with base directory baseDir 
// for template files.
func New(baseDir string, fmap template.FormatterMap) *Manager {
	// Add each built-in formatter unless there's 
	// a user given formatter with same name already.
	if fmap != nil {
		for k, v := range builtins {
			_, present := fmap[k]
			if !present {
				fmap[k] = v
			}
		}
	}

	return &Manager{
		baseDir:   baseDir,
		tStrings:  make(map[string]*Template),
		tFiles:    make(map[string]*Template),
		fmap:      fmap,
		ldelim:    "{",
		rdelim:    "}",
		reloading: false}
}

// Add adds a given template string s to the template manager 
// with the identifier id.
// If any errors occur, returned error will be non-nil. 
func (m *Manager) Add(s string, id string) (*Template, os.Error) {
	return m.add(s, id, false)
}

// AddFile adds a given template file to the template manager.
// If any errors occur, returned error will be non-nil. 
func (m *Manager) AddFile(filename string) (*Template, os.Error) {
	return m.addFile(filename, false)
}

// Removes all templates from the template manager.
// Useful for clearing out cached templates.
// Clear returns true if one or more templates were removed, otherwise false.
func (m *Manager) Clear() bool {
	tlen := len(m.tStrings) + len(m.tFiles)
	m.tStrings = make(map[string]*Template)
	m.tFiles = make(map[string]*Template)
	return tlen > 0
}

// Returns a template with the given identifier or nil if it doesn't exist.
func (m *Manager) Get(s string) *Template {
	return m.tStrings[s]
}

// Returns a template with the given filename or nil if it doesn't exist.
func (m *Manager) GetFile(filename string) *Template {
	return m.tFiles[filename]
}

// MustAdd is like Add, but panics, if template can't be parsed. 
func (m *Manager) MustAdd(s string, id string) *Template {
	t, _ := m.add(s, id, true)
	return t
}

// MustAddFile is like AddFile, but panics, if template can't be parsed. 
func (m *Manager) MustAddFile(filename string) *Template {
	t, _ := m.addFile(filename, true)
	return t
}

// Removes a template with the given identifier from the template manager.
// Useful for clearing out cached templates.
// It's safe to remove a non-existing template.
// Remove returns true if a template was removed, otherwise false.
func (m *Manager) Remove(s string) bool {
	_, present := m.tStrings[s]
	m.tStrings[s] = nil, false
	return present
}

// Removes a template with the given filename from the template manager.
// Useful for clearing out cached templates.
// It's safe to remove a non-existing template.
// Remove returns true if a template was removed, otherwise false.
func (m *Manager) RemoveFile(filename string) bool {
	_, present := m.tFiles[filename]
	m.tFiles[filename] = nil, false
	return present
}

// SetReloading sets the template file reloading mode.
// When reloading mode is enabled, calls to GetFile method will trigger 
// reparsing of the given template file if its modified time has changed.
// Reloading is disabled (false) by default.
func (m *Manager) SetReloading(reloading bool) {
	m.reloading = reloading
}

// SetDelims sets the left and right delimiters for operations 
// in the template for template parsing.
func (m *Manager) SetDelims(left, right string) {
	m.ldelim = left
	m.rdelim = right
}


// Unexported methods

// Add adds a given template string to the template manager.
// This method should not be called directly, but through Add or MustAdd.
// If any errors occur, err will be non-nil. 
func (m *Manager) add(s string, id string, mustParse bool) (t *Template,
err os.Error) {
	tt := template.New(m.fmap)
	tt.SetDelims(m.ldelim, m.rdelim)

	// Parse the template.
	if mustParse {
		err := tt.Parse(s)
		if err != nil {
			panic(err)
		}
	} else {
		err := tt.Parse(s)
		if err != nil {
			return nil, err
		}
	}

	t = &Template{
		m:     m,
		cache: tt}

	// Add template to the manager.
	m.tStrings[id] = t
	return t, nil
}

// AddFile adds a given template file to the template manager.
// This method should not be called directly, but through AddFile, MustAddFile.
// If any errors occur, err will be non-nil. 
func (m *Manager) addFile(filename string, mustParse bool) (t *Template,
err os.Error) {
	var tt *template.Template

	// Parse template file.
	path := path.Join(m.baseDir, filename)
	tt, err = m.Parsett(path, mustParse)
	if err != nil {
		return nil, err
	}

	t = &Template{
		m:     m,
		cache: tt,
		fi: &templateFileInfo{
			filename:  filename,
			mtime:     getMtime(path),
			mustParse: mustParse}}

	// Add template to the manager.
	m.tFiles[filename] = t

	return t, nil
}

// Parsett returns a *template.Template for the given file.
func (m *Manager) Parsett(path string, mustParse bool) (tt *template.Template,
err os.Error) {
	tt = template.New(m.fmap)
	tt.SetDelims(m.ldelim, m.rdelim)

	// Parse template file.
	if mustParse {
		err = tt.ParseFile(path)
		if err != nil {
			panic(err)
		}
	} else {
		err = tt.ParseFile(path)
		if err != nil {
			return nil, err
		}
	}

	return tt, nil
}

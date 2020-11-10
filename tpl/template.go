package tpl

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"path/filepath"
	"strings"
)

// Params -- template parameters (alias to map[string]interface{})
type Params map[string]interface{}

// Templates -- holds a map of compiled templates
type Templates struct {
	// Functions -- list of functions to make available in templates
	//
	// Note that you have to set this before you call Load()
	Functions template.FuncMap

	// if set, calls strings.HasSuffix() for each filename.
	// only templates with a matching suffix will be parsed
	//
	// (make sure to add the leading '.' to the suffix)
	//
	// example:
	// - `[]string{".htm", ".html"}`: only allow .html and .html files
	FileExtensions []string

	// LayoutFiles -- template files that will serve as the base for all others
	// (e.g. a shared layout.html)
	//
	// ignored if empty
	LayoutFiles []string

	// templates -- holds all templates once they've been parsed
	templates map[string]*template.Template
}

// Load -- (re)loads template files in the configured template path
func (t *Templates) Load(templateDir string) error {
	t.templates = make(map[string]*template.Template)

	// load LayoutFiles
	var layout = template.New("layout").Funcs(t.Functions)
	var err error
	for _, layoutPath := range t.LayoutFiles {
		layoutPath = filepath.Join(templateDir, layoutPath) // TODO allow for absolute paths
		if layout, err = layout.ParseFiles(layoutPath); err != nil {
			return fmt.Errorf("failed to load layout '%s': %w", layoutPath, err)
		}
	}

	// recursively load template files
	return t.loadRecursive(layout, templateDir, templateDir)
}

func (t *Templates) loadRecursive(layout *template.Template, dir string, rootPath string) error {
	var children, err = ioutil.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("failed to list template files in '%s': %w", dir, err)
	}

	for _, child := range children {
		var childPath = filepath.Join(dir, child.Name())
		var relPath string // file path relative to rootPath
		if relPath, err = filepath.Rel(rootPath, childPath); err != nil {
			return fmt.Errorf("failed to relativize template path: %w", err)
		}

		if child.IsDir() {
			if err = t.loadRecursive(layout, childPath, rootPath); err != nil {
				return err
			}
		} else if t.filterFilename(relPath) {

			var tpl *template.Template
			if tpl, err = layout.Clone(); err != nil {
				return fmt.Errorf("failed to clone layout: %w", err)
			}

			if _, err = tpl.ParseFiles(childPath); err != nil {
				return fmt.Errorf("failed to parse template '%s': %w", childPath, err)
			}

			// relativize path
			t.templates[relPath] = tpl
		}
	}

	return nil
}

// filterFilename -- returns false if the given file shouldn't be parsed as template
func (t *Templates) filterFilename(path string) bool {
	for _, layoutFile := range t.LayoutFiles {
		if path == layoutFile {
			return false // this is a layout -> don't add
		}
	}

	if len(t.FileExtensions) == 0 {
		return true // no filter
	}

	for _, ext := range t.FileExtensions {
		if strings.HasSuffix(path, ext) {
			return true // valid suffix
		}
	}

	return false
}

// Exists -- returns true if the given template exists
func (t *Templates) Exists(path string) bool {
	var _, ok = t.templates[path]
	return ok
}

// Get -- returns the given template (or nil if not found)
func (t *Templates) Get(path string) *template.Template {
	return t.templates[path]
}

// List -- returns the paths of all loaded templates that have the specified section (e.g. 'message' for emails)
func (t *Templates) List(withSection string) []string {
	var rc = make([]string, 0, len(t.templates))
	for path, tpl := range t.templates {
		if tpl.Lookup(withSection) != nil {
			rc = append(rc, path)
		}
	}
	return rc
}

// MustLoad -- invokes t.Load() and panics on error
func MustLoad(t *Templates, path string) *Templates {
	if err := t.Load(path); err != nil {
		panic(err)
	}
	return t
}

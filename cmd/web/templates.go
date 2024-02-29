package main

import (
	"html/template"
	"path/filepath"

	"snippetbox.flaviogalon.github.io/internal/models"
)

type templateData struct {
	CurrentYear int
	Snippet     *models.Snippet
	Snippets    []*models.Snippet
}

// Return cache of page templates mapped by template name
func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob("./ui/html/pages/*.tmpl.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		// Parse base template file into a template set
		ts, err := template.ParseFiles("./ui/html/base.tmpl.html")
		if err != nil {
			return nil, err
		}

		// Extend the template set with all partials
		ts, err = ts.ParseGlob("./ui/html/partials/*tmpl.html")
		if err != nil {
			return nil, err
		}

		// Finally extend the template set with the page
		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}
	return cache, nil
}

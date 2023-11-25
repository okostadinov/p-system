package main

import (
	"html/template"
	"path/filepath"

	"p-system.okostadinov.net/internal/models"
)

type templateData struct {
	CurrentYear     int
	Patient         *models.Patient
	Patients        []*models.Patient
	Medications     []*models.Medication
	Form            any
	Flash           Flash
	IsAuthenticated bool
	CSRFField       template.HTML
}

// prepares and stores all the html templates upon app initiation
func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob("./ui/html/pages/*.tmpl.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		ts, err := template.ParseFiles("./ui/html/base.tmpl.html")
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob("./ui/html/partials/*.tmpl.html")
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}

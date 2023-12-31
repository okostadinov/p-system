package main

import (
	"html/template"
	"io/fs"
	"path/filepath"

	"p-system.okostadinov.net/internal/models"
	"p-system.okostadinov.net/ui"
)

type templateData struct {
	CurrentYear     int
	Patient         *models.Patient
	Patients        []*models.Patient
	Medications     []*models.Medication
	Form            any
	Flash           Flash
	IsAuthenticated bool
	UserId          int
	CSRFField       template.HTML
}

// prepares and stores all the html templates upon app initiation
func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := fs.Glob(ui.Files, "html/pages/*.tmpl.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		patterns := []string{
			"html/base.tmpl.html",
			"html/partials/*.tmpl.html",
			page,
		}

		ts, err := template.New(name).ParseFS(ui.Files, patterns...)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}

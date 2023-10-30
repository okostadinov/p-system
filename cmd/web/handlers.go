package main

import (
	"html/template"
	"net/http"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	latest, err := app.patients.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	files := []string{
		"./ui/html/base.tmpl.html",
		"./ui/html/partials/nav.tmpl.html",
		"./ui/html/pages/home.tmpl.html",
	}

	t, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}

	err = t.ExecuteTemplate(w, "base", latest)
	if err != nil {
		app.serverError(w, err)
	}
}

func (app *application) create(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("creating a patient"))
}

func (app *application) get(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("getting a patient"))
}

func (app *application) list(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("listing all patients"))
}

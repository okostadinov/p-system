package main

import (
	"html/template"
	"net/http"
)

type patientCreateForm struct {
	UCN         int    `form:"ucn"`
	Name        string `form:"name"`
	PhoneNumber string `form:"phone_number"`
	Height      int    `form:"height"`
	Weight      int    `form:"weight"`
	Medication  string `form:"medication"`
	Note        string `form:"note"`
}

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

func (app *application) patientCreate(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"./ui/html/base.tmpl.html",
		"./ui/html/partials/nav.tmpl.html",
		"./ui/html/pages/create.tmpl.html",
	}

	t, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}

	err = t.ExecuteTemplate(w, "base", nil)
	if err != nil {
		app.serverError(w, err)
	}
}

func (app *application) patientCreatePost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.serverError(w, err)
		return
	}
}

func (app *application) patientList(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("listing all patients"))
}

func (app *application) patientView(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("getting a patient"))
}

func (app *application) patientUpdate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("getting a patient"))
}

func (app *application) patientDelete(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("getting a patient"))
}

func (app *application) medicationList(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("getting a patient"))
}

func (app *application) medicationAdd(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("getting a patient"))
}

func (app *application) medicationDelete(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("getting a patient"))
}

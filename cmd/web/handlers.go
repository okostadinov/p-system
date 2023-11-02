package main

import (
	"database/sql"
	"errors"
	"html/template"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"p-system.okostadinov.net/internal/models"
)

type templateData struct {
	Patient     *models.Patient
	Patients    []*models.Patient
	Medications []*models.Medication
}

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

	data := &templateData{
		Patients: latest,
	}

	err = t.ExecuteTemplate(w, "base", data)
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
	patients, err := app.patients.GetAll()
	if err != nil {
		app.serverError(w, err)
		return
	}

	files := []string{
		"./ui/html/base.tmpl.html",
		"./ui/html/partials/nav.tmpl.html",
		"./ui/html/pages/list.tmpl.html",
	}

	t, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := &templateData{
		Patients: patients,
	}

	err = t.ExecuteTemplate(w, "base", data)
	if err != nil {
		app.serverError(w, err)
	}
}

func (app *application) patientListFiltered(w http.ResponseWriter, r *http.Request) {
	medication := mux.Vars(r)["name"]

	patients, err := app.patients.GetAllByMedication(medication)
	if err != nil {
		app.serverError(w, err)
		return
	}

	files := []string{
		"./ui/html/base.tmpl.html",
		"./ui/html/partials/nav.tmpl.html",
		"./ui/html/pages/list.tmpl.html",
	}

	t, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := &templateData{
		Patients: patients,
	}

	err = t.ExecuteTemplate(w, "base", data)
	if err != nil {
		app.serverError(w, err)
	}
}

func (app *application) patientView(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	patient, err := app.patients.Get(id)
	if err != nil {
		if errors.Is(sql.ErrNoRows, err) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	files := []string{
		"./ui/html/base.tmpl.html",
		"./ui/html/partials/nav.tmpl.html",
		"./ui/html/pages/view.tmpl.html",
	}

	t, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := &templateData{
		Patient: patient,
	}

	err = t.ExecuteTemplate(w, "base", data)
	if err != nil {
		app.serverError(w, err)
	}
}

func (app *application) patientUpdate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("getting a patient"))
}

func (app *application) patientDelete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	err = app.patients.Delete(id)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, "/patients/", http.StatusSeeOther)
}

func (app *application) medicationList(w http.ResponseWriter, r *http.Request) {
	medications, err := app.medications.GetAll()
	if err != nil {
		app.serverError(w, err)
		return
	}

	files := []string{
		"./ui/html/base.tmpl.html",
		"./ui/html/partials/nav.tmpl.html",
		"./ui/html/pages/medications.tmpl.html",
	}

	t, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := &templateData{
		Medications: medications,
	}

	err = t.ExecuteTemplate(w, "base", data)
	if err != nil {
		app.serverError(w, err)
	}
}

func (app *application) medicationAdd(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")

	err := app.medications.Insert(name)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/medications/", http.StatusSeeOther)
}

func (app *application) medicationDelete(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")

	err := app.medications.Delete(name)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, "/medications/", http.StatusSeeOther)
}

package main

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"p-system.okostadinov.net/internal/models"
)

type templateData struct {
	CurrentYear int
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

	data := app.newTemplateData(r)
	data.Patients = latest
	app.render(w, http.StatusOK, "home.tmpl.html", data)
}

func (app *application) patientCreate(w http.ResponseWriter, r *http.Request) {
	medications, err := app.medications.GetAll()
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := app.newTemplateData(r)
	data.Medications = medications
	app.render(w, http.StatusOK, "create.tmpl.html", data)
}

func (app *application) patientCreatePost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	ucn, err := strconv.Atoi(r.PostForm.Get("ucn"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	name := r.PostForm.Get("name")

	number := r.PostForm.Get("phone_number")

	height, err := strconv.Atoi(r.PostForm.Get("height"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	weight, err := strconv.Atoi(r.PostForm.Get("weight"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	medication := r.PostForm.Get("medication")

	note := r.PostForm.Get("note")

	id, err := app.patients.Insert(ucn, name, number, height, weight, medication, note)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/patients/%d", id), http.StatusSeeOther)
}

func (app *application) patientList(w http.ResponseWriter, r *http.Request) {
	patients, err := app.patients.GetAll()
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := app.newTemplateData(r)
	data.Patients = patients
	app.render(w, http.StatusOK, "list.tmpl.html", data)
}

func (app *application) patientListFiltered(w http.ResponseWriter, r *http.Request) {
	medication := mux.Vars(r)["name"]

	patients, err := app.patients.GetAllByMedication(medication)
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := app.newTemplateData(r)
	data.Patients = patients
	app.render(w, http.StatusOK, "list.tmpl.html", data)
}

func (app *application) patientView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
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

	data := app.newTemplateData(r)
	data.Patient = patient
	app.render(w, http.StatusOK, "view.tmpl.html", data)
}

func (app *application) patientUpdate(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	err = r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	ucn, err := strconv.Atoi(r.PostForm.Get("ucn"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	name := r.PostForm.Get("name")

	number := r.PostForm.Get("phone_number")

	height, err := strconv.Atoi(r.PostForm.Get("height"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	weight, err := strconv.Atoi(r.PostForm.Get("weight"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	medication := r.PostForm.Get("medication")

	note := r.PostForm.Get("note")

	approved, err := strconv.ParseBool(r.PostForm.Get("approved"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	firstCont, err := strconv.ParseBool(r.PostForm.Get("first_continuation"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	patient := &models.Patient{
		ID:                id,
		UCN:               ucn,
		Name:              name,
		PhoneNumber:       number,
		Height:            height,
		Weight:            weight,
		Medication:        medication,
		Note:              note,
		Approved:          approved,
		FirstContinuation: firstCont,
	}

	err = app.patients.Update(patient)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/patients/%d", id), http.StatusSeeOther)
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

	data := app.newTemplateData(r)
	data.Medications = medications
	app.render(w, http.StatusOK, "medications.tmpl.html", data)
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

	patients, err := app.patients.GetAllByMedication(name)
	if err != nil {
		app.serverError(w, err)
		return
	}

	if len(patients) > 0 {
		http.Redirect(w, r, "/medications/", http.StatusForbidden)
		// TODO: add message to request context once sessions are added
	}

	err = app.medications.Delete(name)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, "/medications/", http.StatusSeeOther)
}

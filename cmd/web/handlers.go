package main

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"p-system.okostadinov.net/internal/models"
)

type templateData struct {
	CurrentYear int
	Patient     *models.Patient
	Patients    []*models.Patient
	Medications []*models.Medication
	Form        any
}

type patientCreateForm struct {
	UCN         string `schema:"ucn" validate:"required,numeric,len=10"`
	FirstName   string `schema:"first_name" validate:"required,alphaunicode"`
	LastName    string `schema:"last_name" validate:"required,alphaunicode"`
	PhoneNumber string `schema:"phone_number" validate:"required,e164"`
	Height      int    `schema:"height" validate:"required,numeric"`
	Weight      int    `schema:"weight" validate:"required,numeric"`
	Medication  string `schema:"medication" validate:"required"`
	Note        string `schema:"note" validate:"required"`
	FieldErrors map[string]string
}

type patientEditForm struct {
	UCN               string `schema:"ucn" validate:"required,numeric,len=10"`
	FirstName         string `schema:"first_name" validate:"required,alphaunicode"`
	LastName          string `schema:"last_name" validate:"required,alphaunicode"`
	PhoneNumber       string `schema:"phone_number" validate:"required,e164"`
	Height            int    `schema:"height" validate:"required,numeric"`
	Weight            int    `schema:"weight" validate:"required,numeric"`
	Medication        string `schema:"medication" validate:"required"`
	Note              string `schema:"note" validate:"required"`
	Approved          bool   `schema:"approved" validate:"required"`
	FirstContinuation bool   `schema:"first_continuation" validate:"required"`
	FieldErrors       map[string]string
}

type medicationAddForm struct {
	Name        string `schema:"name" validate:"required,alpha"`
	FieldErrors map[string]string
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
	data.Form = &patientCreateForm{}
	app.render(w, http.StatusOK, "create.tmpl.html", data)
}

func (app *application) patientCreatePost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	var form patientCreateForm

	err = app.decoder.Decode(&form, r.PostForm)
	if err != nil {
		app.serverError(w, err)
		return
	}

	form.FieldErrors = map[string]string{}
	err = app.validator.Struct(form)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			form.FieldErrors[err.Field()] = app.fetchTagErrorMessage(err.Tag(), err.Param())
		}
	}

	data := app.newTemplateData(r)
	data.Form = form
	// TODO: replace later with session context
	data.Medications, _ = app.medications.GetAll()

	if len(form.FieldErrors) > 0 {
		app.render(w, http.StatusUnprocessableEntity, "create.tmpl.html", data)
		return
	}

	id, err := app.patients.Insert(form.UCN, form.FirstName, form.LastName, form.PhoneNumber, form.Height, form.Weight, form.Medication, form.Note)
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

	medications, err := app.medications.GetAll()
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := app.newTemplateData(r)
	data.Patient = patient
	data.Medications = medications
	data.Form = &patientEditForm{}
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

	var form patientEditForm
	err = app.decoder.Decode(&form, r.PostForm)
	if err != nil {
		app.serverError(w, err)
		return
	}

	form.FieldErrors = map[string]string{}
	err = app.validator.Struct(form)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			form.FieldErrors[err.Field()] = app.fetchTagErrorMessage(err.Tag(), err.Param())
		}
	}

	data := app.newTemplateData(r)
	data.Form = form
	// TODO: replace later with session context
	data.Patient, _ = app.patients.Get(id)
	data.Medications, _ = app.medications.GetAll()

	if len(form.FieldErrors) > 0 {
		app.render(w, http.StatusUnprocessableEntity, "view.tmpl.html", data)
		return
	}

	err = app.patients.Update(id, form.UCN, form.FirstName, form.LastName, form.PhoneNumber, form.Height, form.Weight, form.Medication, form.Note, form.Approved, form.FirstContinuation)
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
	data.Form = &medicationAddForm{}
	app.render(w, http.StatusOK, "medications.tmpl.html", data)
}

func (app *application) medicationAdd(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	var form medicationAddForm
	err = app.decoder.Decode(&form, r.PostForm)
	if err != nil {
		app.serverError(w, err)
		return
	}

	form.FieldErrors = map[string]string{}
	err = app.validator.Struct(form)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			form.FieldErrors[err.Field()] = app.fetchTagErrorMessage(err.Tag(), err.Param())
		}
	}

	data := app.newTemplateData(r)
	data.Form = form
	data.Medications, _ = app.medications.GetAll()
	if len(form.FieldErrors) > 0 {
		app.render(w, http.StatusUnprocessableEntity, "medications.tmpl.html", data)
		return
	}

	err = app.medications.Insert(form.Name)
	if err != nil {
		app.serverError(w, err)
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

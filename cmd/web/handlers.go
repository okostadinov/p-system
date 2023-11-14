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
	Form        any
	Flash       string
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
	FieldErrors `schema:"-"`
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
	Approved          bool   `schema:"approved" validate:"boolean"`
	FirstContinuation bool   `schema:"first_continuation" validate:"boolean"`
	FieldErrors       `schema:"-"`
}

type medicationAddForm struct {
	Name        string `schema:"name" validate:"required"`
	FieldErrors `schema:"-"`
}

type searchByUCNForm struct {
	UCN string `schema:"q"`
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	latest, err := app.patients.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := app.newTemplateData(w, r)
	data.Patients = latest
	app.render(w, http.StatusOK, "home.tmpl.html", data)
}

func (app *application) patientCreate(w http.ResponseWriter, r *http.Request) {
	medications, err := app.medications.GetAll()
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := app.newTemplateData(w, r)
	data.Medications = medications
	data.Form = &patientCreateForm{}
	app.render(w, http.StatusOK, "create.tmpl.html", data)
}

func (app *application) patientCreatePost(w http.ResponseWriter, r *http.Request) {
	var form patientCreateForm
	err := app.decodeForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	if ok, errors := app.validateForm(form); !ok {
		medications, err := app.medications.GetAll()
		if err != nil {
			app.serverError(w, err)
			return
		}

		data := app.newTemplateData(w, r)
		data.Medications = medications
		form.FieldErrors = errors
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "create.tmpl.html", data)
		return
	}

	id, err := app.patients.Insert(form.UCN, form.FirstName, form.LastName, form.PhoneNumber, form.Height, form.Weight, form.Medication, form.Note)
	if err != nil {
		app.serverError(w, err)
		return
	}

	err = app.setFlash(w, r, "Пациентът бе добавен успешно")
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

	data := app.newTemplateData(w, r)
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

	data := app.newTemplateData(w, r)
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

	data := app.newTemplateData(w, r)
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

	var form patientEditForm
	err = app.decodeForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	if ok, errors := app.validateForm(form); !ok {
		medications, err := app.medications.GetAll()
		if err != nil {
			app.serverError(w, err)
			return
		}

		patient, err := app.patients.Get(id)
		if err != nil {
			app.serverError(w, err)
			return
		}

		data := app.newTemplateData(w, r)
		data.Medications = medications
		form.FieldErrors = errors
		data.Form = form
		data.Patient = patient
		app.render(w, http.StatusUnprocessableEntity, "view.tmpl.html", data)
		return
	}

	err = app.patients.Update(id, form.UCN, form.FirstName, form.LastName, form.PhoneNumber, form.Height, form.Weight, form.Medication, form.Note, form.Approved, form.FirstContinuation)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.setFlash(w, r, "Данните бяха обновени успешно")
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

	app.setFlash(w, r, "Пациентът бе изтрит успешно")
	http.Redirect(w, r, "/patients/", http.StatusSeeOther)
}

func (app *application) patientSearchByUCN(w http.ResponseWriter, r *http.Request) {
	var form searchByUCNForm
	err := app.decodeForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	patient, err := app.patients.GetByUCN(form.UCN)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			app.setFlash(w, r, "Не съществува пациент с такъв ЕГН")
			http.Redirect(w, r, "/", http.StatusSeeOther)
		} else {
			app.serverError(w, err)
		}
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/patients/%d", patient.ID), http.StatusSeeOther)
}

func (app *application) medicationList(w http.ResponseWriter, r *http.Request) {
	medications, err := app.medications.GetAll()
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := app.newTemplateData(w, r)
	data.Medications = medications
	data.Form = &medicationAddForm{}
	app.render(w, http.StatusOK, "medications.tmpl.html", data)
}

func (app *application) medicationAdd(w http.ResponseWriter, r *http.Request) {
	var form medicationAddForm
	err := app.decodeForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	if ok, errors := app.validateForm(form); !ok {
		medications, err := app.medications.GetAll()
		if err != nil {
			app.serverError(w, err)
			return
		}

		data := app.newTemplateData(w, r)
		form.FieldErrors = errors
		data.Form = form
		data.Medications = medications
		app.render(w, http.StatusUnprocessableEntity, "medications.tmpl.html", data)
		return
	}

	err = app.medications.Insert(form.Name)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.setFlash(w, r, "Медикаментът бе добавен успешно")
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
		app.setFlash(w, r, "Медикаментът не може да бъде изтрит, поради записани с него пациенти")
		http.Redirect(w, r, "/medications/", http.StatusSeeOther)
		return
	}

	err = app.medications.Delete(name)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.setFlash(w, r, "Медикаментът бе изтрит успешно")
	http.Redirect(w, r, "/medications/", http.StatusSeeOther)
}

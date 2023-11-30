package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"p-system.okostadinov.net/internal/models"
	"p-system.okostadinov.net/internal/validator"
)

type patientForm struct {
	UCN                  string `schema:"ucn" validate:"required,numeric,len=10"`
	FirstName            string `schema:"first_name" validate:"required,alphaunicode"`
	LastName             string `schema:"last_name" validate:"required,alphaunicode"`
	PhoneNumber          string `schema:"phone_number" validate:"required,e164"`
	Height               int    `schema:"height" validate:"required,numeric"`
	Weight               int    `schema:"weight" validate:"required,numeric"`
	Medication           string `schema:"medication" validate:"required"`
	Note                 string `schema:"note" validate:"required"`
	Approved             bool   `schema:"approved" validate:"boolean"`
	FirstContinuation    bool   `schema:"first_continuation" validate:"boolean"`
	validator.FormErrors `schema:"-"`
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
	data.Form = &patientForm{}
	app.render(w, http.StatusOK, "create.tmpl.html", data)
}

func (app *application) patientCreatePost(w http.ResponseWriter, r *http.Request) {
	var form patientForm
	err := app.decodeForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	if !app.validator.ValidateForm(form) {
		medications, err := app.medications.GetAll()
		if err != nil {
			app.serverError(w, err)
			return
		}

		data := app.newTemplateData(w, r)
		data.Medications = medications
		form.FormErrors = app.validator.FormErrors
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "create.tmpl.html", data)
		return
	}

	userId := app.getUserIdFromContext(w, r)
	if userId == 0 {
		app.serverError(w, err)
		return
	}

	id, err := app.patients.Insert(form.UCN, form.FirstName, form.LastName, form.PhoneNumber, form.Height, form.Weight, form.Medication, form.Note, userId)
	if err != nil {
		app.serverError(w, err)
		return
	}

	err = app.setFlash(w, r, "Patient successfully added!", FlashTypeSuccess)
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
		if errors.Is(err, models.ErrNoRecord) {
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
	data.Form = &patientForm{}
	app.render(w, http.StatusOK, "view.tmpl.html", data)
}

func (app *application) patientUpdate(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	var form patientForm
	err = app.decodeForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	if !app.validator.ValidateForm(form) {
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
		form.FormErrors = app.validator.FormErrors
		data.Form = form
		data.Patient = patient
		app.render(w, http.StatusUnprocessableEntity, "view.tmpl.html", data)
		return
	}

	err = app.patients.Update(id, form.UCN, form.FirstName, form.LastName, form.PhoneNumber, form.Height, form.Weight, form.Medication, form.Note, form.Approved, form.FirstContinuation, app.getUserIdFromContext(w, r))
	if err != nil {
		if errors.Is(err, models.ErrUnauthorizedAction) {
			err = app.setFlash(w, r, "Unauthorized action - cannot modify patient!", FlashTypeDanger)
			if err != nil {
				app.serverError(w, err)
				return
			}
			http.Redirect(w, r, fmt.Sprintf("/patients/%d", id), http.StatusSeeOther)
		} else {
			app.serverError(w, err)
		}
		return
	}

	err = app.setFlash(w, r, "Patient successfully updated!", FlashTypeSuccess)
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

	err = app.patients.Delete(id, app.getUserIdFromContext(w, r))
	if err != nil {
		if errors.Is(err, models.ErrUnauthorizedAction) {
			err = app.setFlash(w, r, "Unauthorized action - cannot delete patient!", FlashTypeDanger)
			if err != nil {
				app.serverError(w, err)
				return
			}
			http.Redirect(w, r, "/patients/", http.StatusSeeOther)
		} else {
			app.serverError(w, err)
		}
		return
	}

	err = app.setFlash(w, r, "Patient successfully deleted!", FlashTypeSuccess)
	if err != nil {
		app.serverError(w, err)
		return
	}
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
		if errors.Is(err, models.ErrNoRecord) {
			err = app.setFlash(w, r, "No patients exists with this UCN.", FlashTypeWarning)
			if err != nil {
				app.serverError(w, err)
				return
			}
			http.Redirect(w, r, r.Header.Get("Referer"), http.StatusSeeOther)
		} else {
			app.serverError(w, err)
		}
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/patients/%d", patient.ID), http.StatusSeeOther)
}

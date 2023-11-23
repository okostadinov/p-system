package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
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

type userSignupForm struct {
	Name            string `schema:"name" validate:"required"`
	Email           string `schema:"email" validate:"required,email"`
	Password        string `schema:"password" validate:"required,password"`
	ConfirmPassword string `schema:"confirm_password" validate:"required,password,eqfield=Password"`
	FieldErrors     `schema:"-"`
}

type userLoginForm struct {
	Email       string `schema:"email" validate:"required,email"`
	Password    string `schema:"password" validate:"required,password"`
	FieldErrors `schema:"-"`
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

	err = app.setFlash(w, r, "Patient successfully added!", "success")
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

	err = app.setFlash(w, r, "Patient successfully updated!", "success")
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

	err = app.setFlash(w, r, "Patient successfully deleted!", "success")
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
			err = app.setFlash(w, r, "No patients exists with this UCN.", "warning")
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

	err = app.setFlash(w, r, "Medication successfully added!", "success")
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
		err = app.setFlash(w, r, "Medication cannot be deleted due to registed patients.", "warning")
		if err != nil {
			app.serverError(w, err)
			return
		}
		http.Redirect(w, r, "/medications/", http.StatusSeeOther)
		return
	}

	err = app.medications.Delete(name)
	if err != nil {
		app.serverError(w, err)
		return
	}

	err = app.setFlash(w, r, "Medication successfully deleted!", "success")
	if err != nil {
		app.serverError(w, err)
		return
	}
	http.Redirect(w, r, "/medications/", http.StatusSeeOther)
}

func (app *application) userSignup(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(w, r)
	data.Form = &userSignupForm{}
	app.render(w, http.StatusOK, "signup.tmpl.html", data)
}

func (app *application) userSignupPost(w http.ResponseWriter, r *http.Request) {
	var form userSignupForm
	err := app.decodeForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	if ok, errors := app.validateForm(form); !ok {
		data := app.newTemplateData(w, r)
		form.FieldErrors = errors
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "signup.tmpl.html", data)
		return
	}

	err = app.users.Insert(form.Name, form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			data := app.newTemplateData(w, r)
			form.FieldErrors = make(FieldErrors)
			form.FieldErrors["email"] = "email address already in use"
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, "signup.tmpl.html", data)
		} else {
			app.serverError(w, err)
		}
		return
	}

	err = app.setFlash(w, r, "Registration successful! You may now log in.", "success")
	if err != nil {
		app.serverError(w, err)
		return
	}
	http.Redirect(w, r, "/users/login", http.StatusSeeOther)
}

func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(w, r)
	data.Form = &userLoginForm{}
	app.render(w, http.StatusOK, "login.tmpl.html", data)
}

func (app *application) userLoginPost(w http.ResponseWriter, r *http.Request) {
	var form userLoginForm
	err := app.decodeForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	if ok, errors := app.validateForm(form); !ok {
		data := app.newTemplateData(w, r)
		form.FieldErrors = errors
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "login.tmpl.html", data)
		return
	}

	id, err := app.users.Authenticate(form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			err = app.setFlash(w, r, "Invalid email address or password.", "danger")
			if err != nil {
				app.serverError(w, err)
				return
			}
			http.Redirect(w, r, "/users/login", http.StatusSeeOther)
		} else {
			app.serverError(w, err)
		}
		return
	}

	session, err := app.store.Get(r, "session")
	if err != nil {
		app.serverError(w, err)
		return
	}
	session.Values["authenticatedUserID"] = id

	err = app.setFlash(w, r, "Logged in successfully!", "success")
	if err != nil {
		app.serverError(w, err)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) userLogout(w http.ResponseWriter, r *http.Request) {
	session, err := app.store.Get(r, "session")
	if err != nil {
		app.serverError(w, err)
		return
	}
	delete(session.Values, "authenticatedUserID")

	err = app.setFlash(w, r, "Logged out successfully!", "success")
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

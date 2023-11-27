package main

import (
	"net/http"

	"p-system.okostadinov.net/internal/validator"
)

type medicationAddForm struct {
	Name                 string `schema:"name" validate:"required"`
	validator.FormErrors `schema:"-"`
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

	if !app.validator.ValidateForm(form) {
		medications, err := app.medications.GetAll()
		if err != nil {
			app.serverError(w, err)
			return
		}

		data := app.newTemplateData(w, r)
		form.FormErrors = app.validator.FormErrors
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

	err = app.setFlash(w, r, "Medication successfully added!", FlashTypeSuccess)
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
		err = app.setFlash(w, r, "Medication cannot be deleted due to registed patients.", FlashTypeWarning)
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

	err = app.setFlash(w, r, "Medication successfully deleted!", FlashTypeSuccess)
	if err != nil {
		app.serverError(w, err)
		return
	}
	http.Redirect(w, r, "/medications/", http.StatusSeeOther)
}

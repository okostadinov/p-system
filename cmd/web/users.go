package main

import (
	"errors"
	"net/http"

	"p-system.okostadinov.net/internal/models"
	"p-system.okostadinov.net/internal/validator"
)

type userSignupForm struct {
	Name                 string `schema:"name" validate:"required"`
	Email                string `schema:"email" validate:"required,email"`
	Password             string `schema:"password" validate:"required,password"`
	ConfirmPassword      string `schema:"confirm_password" validate:"required,password,eqfield=Password"`
	validator.FormErrors `schema:"-"`
}

type userLoginForm struct {
	Email                string `schema:"email" validate:"required,email"`
	Password             string `schema:"password" validate:"required,password"`
	validator.FormErrors `schema:"-"`
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

	if !app.validator.ValidateForm(form) {
		data := app.newTemplateData(w, r)
		form.FormErrors = app.validator.FormErrors
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "signup.tmpl.html", data)
		return
	}

	err = app.users.Insert(form.Name, form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			data := app.newTemplateData(w, r)
			app.validator.FormErrors["email"] = "email address already in use"
			form.FormErrors = app.validator.FormErrors
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

	if !app.validator.ValidateForm(form) {
		data := app.newTemplateData(w, r)
		form.FormErrors = app.validator.FormErrors
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

	session.Values["userID"] = id
	err = session.Save(r, w)
	if err != nil {
		app.serverError(w, err)
		return
	}

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

	delete(session.Values, "userID")
	err = session.Save(r, w)
	if err != nil {
		app.serverError(w, err)
		return
	}

	err = app.setFlash(w, r, "Logged out successfully!", "success")
	if err != nil {
		app.serverError(w, err)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/go-playground/validator/v10"
)

type FieldErrors map[string]string

// outputs the error to the client, as well as logging it locally for debugging
func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// outputs the status code as well as the text associated with it to the client
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// for consistency, instead of having to call 'http.NotFound()' exclusively for 404 errors
func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

// retrieves and executes a particular html template from the app's template cache
func (app *application) render(w http.ResponseWriter, status int, page string, data *templateData) {
	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		app.serverError(w, err)
		return
	}

	buf := new(bytes.Buffer)

	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(w, err)
		return
	}

	w.WriteHeader(status)

	buf.WriteTo(w)
}

// prepares a template data struct with common dynamic data
func (app *application) newTemplateData(w http.ResponseWriter, r *http.Request) *templateData {
	return &templateData{
		CurrentYear: time.Now().Year(),
		Flash:       app.popFlash(w, r),
	}
}

// decodes the request into a form struct
func (app *application) decodeForm(r *http.Request, form interface{}) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}

	err = app.decoder.Decode(form, r.PostForm)
	if err != nil {
		return err
	}

	return nil
}

// parses the form and returns whether the validation was successful or not, together with the a map of errors
func (app *application) validateForm(form interface{}) (bool, FieldErrors) {
	err := app.validator.Struct(form)

	if err != nil {
		errors := make(FieldErrors)

		for _, err := range err.(validator.ValidationErrors) {
			errors[err.Field()] = app.fetchTagErrorMessage(err.Tag(), err.Param())
		}

		return false, errors
	}

	return true, nil
}

// gets a translated form field error message based on the error
func (app *application) fetchTagErrorMessage(tag, param string) string {
	switch tag {
	case "required":
		return "задължително поле"
	case "numeric":
		return "грешен формат (допустими стойности - цифри)"
	case "len":
		return fmt.Sprintf("невалидно количество символи (нужни - %v)", param)
	case "alphaunicode":
		return "грешен формат (допустими стойности - букви)"
	case "e164":
		return "грешен формат номер (пр. +359123456789)"
	default:
		return ""
	}
}

// adds a flash message to the current session
func (app *application) setFlash(w http.ResponseWriter, r *http.Request, msg string) error {
	session, err := app.store.Get(r, "session")
	if err != nil {
		return err
	}

	session.AddFlash(msg, "message")
	err = session.Save(r, w)
	if err != nil {
		return err
	}

	return nil
}

// pops the flash message from the current session and returns it
func (app *application) popFlash(w http.ResponseWriter, r *http.Request) string {
	session, err := app.store.Get(r, "session")
	if err != nil {
		return ""
	}

	var flashMsg string
	if flashes := session.Flashes("message"); len(flashes) > 0 {
		flashMsg = flashes[0].(string)
	}

	err = session.Save(r, w)
	if err != nil {
		return ""
	}

	return flashMsg
}

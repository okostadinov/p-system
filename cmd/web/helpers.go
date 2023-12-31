package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/gorilla/csrf"
)

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
		CurrentYear:     time.Now().Year(),
		Flash:           app.popFlash(w, r),
		IsAuthenticated: app.isAuthenticated(w, r),
		UserId:          app.getUserIdFromContext(w, r),
		CSRFField:       csrf.TemplateField(r),
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

// checks whether a user is logged in based on the request context
func (app *application) isAuthenticated(w http.ResponseWriter, r *http.Request) bool {
	isAuthenticated, ok := r.Context().Value(isAuthenticatedContextKey).(bool)
	if !ok {
		return false
	}
	return isAuthenticated
}

// retrieves the current authenticated user's ID from the current session
func (app *application) getUserId(w http.ResponseWriter, r *http.Request) int {
	session, _ := app.store.Get(r, "session")
	userID, ok := session.Values["userID"].(int)
	if !ok {
		return 0
	}
	return userID
}

// fetches the current authenticated user's ID from the request context
func (app *application) getUserIdFromContext(w http.ResponseWriter, r *http.Request) int {
	userId, ok := r.Context().Value(userIdContextKey).(int)
	if !ok {
		return 0
	}
	return userId
}

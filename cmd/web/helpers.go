package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
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

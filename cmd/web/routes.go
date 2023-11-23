package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

// registers the routes to a mux assigned to the server
func (app *application) routes() http.Handler {
	mux := mux.NewRouter()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fileServer))

	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	protected := standard.Append(app.requireAuthentication)

	mux.Handle("/", standard.ThenFunc(app.home)).Methods("GET")

	patientsRouter := mux.PathPrefix("/patients").Subrouter()
	patientsRouter.Handle("/", protected.ThenFunc(app.patientList)).Methods("GET")
	patientsRouter.Handle("/medication/{name}", protected.ThenFunc(app.patientListFiltered)).Methods("GET")
	patientsRouter.Handle("/create", protected.ThenFunc(app.patientCreate)).Methods("GET")
	patientsRouter.Handle("/create", protected.ThenFunc(app.patientCreatePost)).Methods("POST")
	patientsRouter.Handle("/delete", protected.ThenFunc(app.patientDelete)).Methods("POST")
	patientsRouter.Handle("/{id:[0-9]+}", protected.ThenFunc(app.patientView)).Methods("GET")
	patientsRouter.Handle("/{id:[0-9]+}", protected.ThenFunc(app.patientUpdate)).Methods("POST")
	patientsRouter.Handle("/search", protected.ThenFunc(app.patientSearchByUCN)).Methods("POST")

	medicationsRouter := mux.PathPrefix("/medications").Subrouter()
	medicationsRouter.Handle("/", protected.ThenFunc(app.medicationList)).Methods("GET")
	medicationsRouter.Handle("/", protected.ThenFunc(app.medicationAdd)).Methods("POST")
	medicationsRouter.Handle("/delete", protected.ThenFunc(app.medicationDelete)).Methods("POST")

	userRouter := mux.PathPrefix("/users").Subrouter()
	userRouter.Handle("/signup", standard.ThenFunc(app.userSignup)).Methods("GET")
	userRouter.Handle("/signup", standard.ThenFunc(app.userSignupPost)).Methods("POST")
	userRouter.Handle("/login", standard.ThenFunc(app.userLogin)).Methods("GET")
	userRouter.Handle("/login", standard.ThenFunc(app.userLoginPost)).Methods("POST")
	userRouter.Handle("/logout", protected.ThenFunc(app.userLogout)).Methods("POST")

	return mux
}

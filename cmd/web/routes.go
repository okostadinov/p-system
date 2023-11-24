package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

// registers the routes to a mux assigned to the server
func (app *application) routes() http.Handler {
	mux := mux.NewRouter()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fileServer))

	mux.Use(app.recoverPanic, app.logRequest, secureHeaders)

	mux.HandleFunc("/", app.home).Methods("GET")

	patientsRouter := mux.PathPrefix("/patients").Subrouter()
	patientsRouter.Use(app.requireAuthentication)
	patientsRouter.HandleFunc("/", app.patientList).Methods("GET")
	patientsRouter.HandleFunc("/medication/{name}", app.patientListFiltered).Methods("GET")
	patientsRouter.HandleFunc("/create", app.patientCreate).Methods("GET")
	patientsRouter.HandleFunc("/create", app.patientCreatePost).Methods("POST")
	patientsRouter.HandleFunc("/delete", app.patientDelete).Methods("POST")
	patientsRouter.HandleFunc("/{id:[0-9]+}", app.patientView).Methods("GET")
	patientsRouter.HandleFunc("/{id:[0-9]+}", app.patientUpdate).Methods("POST")
	patientsRouter.HandleFunc("/search", app.patientSearchByUCN).Methods("POST")

	medicationsRouter := mux.PathPrefix("/medications").Subrouter()
	medicationsRouter.Use(app.requireAuthentication)
	medicationsRouter.HandleFunc("/", app.medicationList).Methods("GET")
	medicationsRouter.HandleFunc("/", app.medicationAdd).Methods("POST")
	medicationsRouter.HandleFunc("/delete", app.medicationDelete).Methods("POST")

	userRouter := mux.PathPrefix("/users").Subrouter()
	userRouter.HandleFunc("/signup", app.userSignup).Methods("GET")
	userRouter.HandleFunc("/signup", app.userSignupPost).Methods("POST")
	userRouter.HandleFunc("/login", app.userLogin).Methods("GET")
	userRouter.HandleFunc("/login", app.userLoginPost).Methods("POST")

	userRouterProtected := mux.PathPrefix("/users").Subrouter()
	userRouterProtected.Use(app.requireAuthentication)
	userRouterProtected.HandleFunc("/logout", app.userLogout).Methods("POST")

	return mux
}

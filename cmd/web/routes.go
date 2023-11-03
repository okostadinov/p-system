package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

// registers the routes to a mux assigned to the server
func (app *application) routes() *mux.Router {
	mux := mux.NewRouter()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fileServer))

	mux.HandleFunc("/", app.home).Methods("GET")

	patientsRouter := mux.PathPrefix("/patients").Subrouter()
	patientsRouter.HandleFunc("/", app.patientList).Methods("GET")
	patientsRouter.HandleFunc("/medication/{name}", app.patientListFiltered).Methods("GET")
	patientsRouter.HandleFunc("/create", app.patientCreate).Methods("GET")
	patientsRouter.HandleFunc("/create", app.patientCreatePost).Methods("POST")
	patientsRouter.HandleFunc("/delete", app.patientDelete).Methods("POST")
	patientsRouter.HandleFunc("/{id:[0-9]+}", app.patientView).Methods("GET")
	patientsRouter.HandleFunc("/{id:[0-9]+}", app.patientUpdate).Methods("POST")

	medicationsRouter := mux.PathPrefix("/medications").Subrouter()
	medicationsRouter.HandleFunc("/", app.medicationList).Methods("GET")
	medicationsRouter.HandleFunc("/", app.medicationAdd).Methods("POST")
	medicationsRouter.HandleFunc("/delete", app.medicationDelete).Methods("POST")
	return mux
}

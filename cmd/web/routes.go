package main

import (
	"github.com/gorilla/mux"
)

// registers the routes to a mux assigned to the server
func (app *application) routes() *mux.Router {
	mux := mux.NewRouter()

	mux.HandleFunc("/", app.home).Methods("GET")

	patientsRouter := mux.PathPrefix("/patients").Subrouter()
	patientsRouter.HandleFunc("/", app.patientList).Methods("GET")
	patientsRouter.HandleFunc("/create", app.patientCreate).Methods("GET")
	patientsRouter.HandleFunc("/create", app.patientCreatePost).Methods("POST")
	patientsRouter.HandleFunc("/{id}", app.patientView).Methods("GET")
	patientsRouter.HandleFunc("/{id}", app.patientUpdate).Methods("PUT")
	patientsRouter.HandleFunc("/{id}", app.patientDelete).Methods("DELETE")

	medicationsRouter := mux.PathPrefix("/medications").Subrouter()
	medicationsRouter.HandleFunc("/", app.medicationList).Methods("GET")
	medicationsRouter.HandleFunc("/", app.medicationAdd).Methods("POST")
	medicationsRouter.HandleFunc("/{id}", app.medicationDelete).Methods("DELETE")
	return mux
}

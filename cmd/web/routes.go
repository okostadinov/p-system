package main

import "net/http"

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/list", app.list)
	mux.HandleFunc("/get", app.get)
	mux.HandleFunc("/create", app.create)

	return mux
}

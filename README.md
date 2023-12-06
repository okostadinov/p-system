# P-System

### Description

P-System is a webapp used to manage patient and medications data using a local SQLite database.

### Features

* CRUD operations for patients and medications
* filter patients by medication
* search for a specific patient based on national ID
* dynamic html
* form validations
* session flash messages

### Setup

* requires `gcc` and the env variable `CGO_ENABLED` to be set
* clone the repo
* `cd` into the project folder
* to run the project `go run ./cmd/web`
* to build an executable `go build ./cmd/web`

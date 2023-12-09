# P-System

#### Note

For a simpler local implementation aimed at working on mobile devices as well, see branch **sqlite-noauth**.

### Description

P-System is a webapp developed with GO connecting to a MySQL database used to store patient and medication data.
It implements a web fronted in order to access the endpoints. Basic authentication is used for authorization
of its users granting access to the system's functionalities. GET requests are available to all users upon
being authenticated, while POST requests are limited only to own patients and medications. For example Mark
can see Emily's patients, but he cannot modify them.

### Features

* CRUD operations for patients and medications
* filtering of patients based on medication
* filtering only own created patients
* looking up patients by UCN (ID)
* dynamic html templating
* form validations
* sessions (incl flash messages)
* authentication & authorization
    * all authenticated users gain access to viewing all patients and medications
    * users may update and delete only own created patients and medications
* static files and template embedding for a self-sufficient binary

### Setup

* clone the repo and `cd` to it
* you will need a TLS certificate:
    * `mkdir tls` and put the `cert.pem` and `key.pem` files into the tls folder
    * for local development `cd tls` and `go run <path-to-GO-stdlib>/src/crypto/tls/generate_cert.go --rsa-bits=2048 --host=localhost`
* run the database setup script (if via terminal: `sudo mysql -u root -p < ./scripts/setup.sql`)
* to start up the project `go run ./cmd/web`
* to build an executable `go build ./cmd/web`
* to see flags usage, append `-h`/`--help` to run command

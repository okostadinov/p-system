# P-System

### Description

P-System is a webapp developed with GO connecting to a MySQL database used to store patient and medication data.
It implements a web fronted in order to access the endpoints. Basic authentication is used for authorization
of its users granting access to the system's functionalities. The patients and medications tables are shared.

### Features

* CRUD operations for patients and medications
* filtering of patients based on medication
* looking up patients by UCN (national ID)
* dynamic html templating
* form validations
* sessions (incl flash messages)
* authentication & authorization

### Setup

* clone the repo and `cd` to it
* you will need a TLS certificate:
    * `mkdir tls` and put the `cert.pem` and `key.pem` files into the tls folder
    * for local dev `cd tls` and `go run <path-to-GO-stdlib>/src/crypto/tls/generate_cert.go --rsa-bits=2048 --host=localhost`
* run the db setup script (if via terminal: `sudo mysql -u root -p < ./scripts/setup.sql`)
* to start up the project `go run ./cmd/web`
* to build an executable `go build ./cmd/web`
* to see flags usage, append `-h`/`--help` to run command

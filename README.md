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
* `mkdir tls && cd tls`
* generate a self-signed TLS cert (e.g. `go run <path-to-GO-stdlib>/src/crypto/tls/generate_cert.go --rsa-bits=2048 --host=localhost`)
* run `setup.sql` (via terminal: `sudo mysql -u root -p < setup.sql`)
* to run the project `go run ./cmd/web`
* to build an executable `go build ./cmd/web`
* to see flags usage, append `-h`/`--help` to run command

# P-System

### Description

P-System is a webapp used to manage patient and medications data.
This branch is intended for use as a locally hosted app which is why an SQLite database was chosen.

### Features

* CRUD operations for patients and medications
* filter patients by medication
* search for a specific patient based on national ID
* dynamic html
* form validations
* session flash messages
* embedded static files and templates for a self-sufficient binary
* auto open a browser window and navigate to the webapp
* shutdown server upon closing the browser

### Setup

* requires `gcc` and the env variable `CGO_ENABLED=1`
* clone the repo
* `cd` into the project folder
* to run the project `go run ./cmd/web`
* to build an executable `go build ./cmd/web`

#### Building an executable for Android
* install Android NDK (e.g. via Android Studio -> SDK Manager -> SDK Tools -> NDK)
  * as of writing the latest release is **26.1.10909125**
* you will need to reference the path to the compiler (in my case: `$HOME/Android/Sdk/ndk/26.1.10909125`)
  * for ease of use you might want to add it as a variable `echo "export NDK_ROOT=<path-to-ndk>" >> ~/.bashrc`
* from the project root run `GOOS=android GOARCH=arm64 CGO_ENABLED=1 CC=$NDK_ROOT/toolchains/llvm/prebuilt/linux-x86_64/bin/aarch64-linux-android34-clang go build ./cmd/web/`
* transfer the **web** executable to the Android device
* in order to run it, you will need a terminal emulator (e.g. Termux via FDroid)
  * run `termux-storage-setup` in order to grant storage access permissions
* find the binary and attempt to run it `./web`
  * in case of insufficient permissions, run `chmod +x ./web`

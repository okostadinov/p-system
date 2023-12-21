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

### Setup

* requires `gcc` and the env variable `CGO_ENABLED=1`
* clone the repo
* `cd` into the project folder
* to run the project without building `make run`
* `make build` will create an executable compatible with the current OS
* `make compile` will create executables for linux, windows and darwin 

#### Building an executable for Android
* install Android NDK (e.g. via Android Studio -> SDK Manager -> SDK Tools -> NDK)
  * as of writing the latest release is **26.1.10909125**
* you will need to reference the path (in my case: `$HOME/Android/Sdk/ndk/26.1.10909125`)
  * for ease of use you might want to add it as a bash variable `echo "export ANDROID_NDK_HOME=<path-to-ndk>" >> ~/.bashrc`
* from the project root run `make android` in order to build the executable at `dist/p_system-android`
* in order to run it on the device, you will need a terminal emulator (e.g. Termux via FDroid)
  * run `termux-setup-storage` in order to grant storage access permissions
  * locate the binary and grant it execute permissions via `chmod +x <<path-to-binary>>`
* now you should be able to run it from the terminal

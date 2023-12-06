package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"reflect"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/michaeljs1990/sqlitestore"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/schema"
	"github.com/gorilla/sessions"
	"p-system.okostadinov.net/internal/models"
)

type application struct {
	infoLog       *log.Logger
	errorLog      *log.Logger
	medications   *models.MedicationModel
	patients      *models.PatientModel
	templateCache map[string]*template.Template
	decoder       *schema.Decoder
	validator     *validator.Validate
	store         *sqlitestore.SqliteStore
}

var validate *validator.Validate

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	dbpath := flag.String("dbpath", "./p_system.db", "SQLite database file path")
	storeKey := flag.String("storekey", "secretkey", "SQLite session store key")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(*dbpath)
	if err != nil {
		errorLog.Fatal(err)
	}

	defer db.Close()

	err = models.SetupDB(db)
	if err != nil {
		errorLog.Fatal(err)
	}

	templateCache, err := newTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}

	store, err := sqlitestore.NewSqliteStoreFromConnection(db, "sessions", "/", 3600, []byte(*storeKey))
	if err != nil {
		errorLog.Fatal(err)
	}

	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   864000 * 7,
		HttpOnly: true,
		Secure:   true,
	}

	defer store.Close()

	app := &application{
		infoLog:       infoLog,
		errorLog:      errorLog,
		medications:   &models.MedicationModel{DB: db},
		patients:      &models.PatientModel{DB: db},
		templateCache: templateCache,
		decoder:       schema.NewDecoder(),
		validator:     setupValidator(),
		store:         store,
	}

	srv := &http.Server{
		Addr:         *addr,
		ErrorLog:     errorLog,
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	errorLog.Fatal(srv.ListenAndServe())
}

func openDB(dbpath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbpath)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func setupValidator() *validator.Validate {
	validate = validator.New(validator.WithRequiredStructEnabled())
	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		return field.Tag.Get("schema")
	})

	return validate
}

package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"reflect"

	_ "github.com/go-sql-driver/mysql"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/schema"
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
}

var validate *validator.Validate

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	dsn := flag.String("dsn", "admin:admin@/p_system?parseTime=true", "MySQL data source name")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	defer db.Close()

	templateCache, err := newTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}

	decoder := schema.NewDecoder()
	validate = validator.New(validator.WithRequiredStructEnabled())
	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		return field.Tag.Get("schema")
	})

	app := &application{
		infoLog:       infoLog,
		errorLog:      errorLog,
		medications:   &models.MedicationModel{DB: db},
		patients:      &models.PatientModel{DB: db},
		templateCache: templateCache,
		decoder:       decoder,
		validator:     validate,
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	errorLog.Fatal(srv.ListenAndServe())
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

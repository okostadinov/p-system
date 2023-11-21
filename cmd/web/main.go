package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"reflect"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/srinathgs/mysqlstore"

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
	users         *models.UserModel
	templateCache map[string]*template.Template
	decoder       *schema.Decoder
	validator     *validator.Validate
	store         *mysqlstore.MySQLStore
}

var validate *validator.Validate

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	dsn := flag.String("dsn", "admin:admin@/p_system?parseTime=true&loc=Local", "MySQL data source name")
	storeKey := flag.String("storekey", "secretkey", "MySQL session store key")
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

	store, err := mysqlstore.NewMySQLStoreFromConnection(db, "sessions", "/", 3600, []byte(*storeKey))
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
		users:         &models.UserModel{DB: db},
		templateCache: templateCache,
		decoder:       schema.NewDecoder(),
		validator:     setupValidator(),
		store:         store,
	}

	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	srv := &http.Server{
		Addr:         *addr,
		ErrorLog:     errorLog,
		Handler:      app.routes(),
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	errorLog.Fatal(srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem"))
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

func setupValidator() *validator.Validate {
	validate = validator.New(validator.WithRequiredStructEnabled())
	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		return field.Tag.Get("schema")
	})

	validate.RegisterValidation("password", passwordValidate)
	return validate
}

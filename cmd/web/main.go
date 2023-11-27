package main

import (
	"crypto/tls"
	"database/sql"
	"encoding/gob"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/srinathgs/mysqlstore"

	"github.com/gorilla/schema"
	"github.com/gorilla/sessions"
	"p-system.okostadinov.net/internal/models"
	"p-system.okostadinov.net/internal/validator"
)

type application struct {
	infoLog       *log.Logger
	errorLog      *log.Logger
	medications   *models.MedicationModel
	patients      *models.PatientModel
	users         *models.UserModel
	templateCache map[string]*template.Template
	decoder       *schema.Decoder
	validator     *validator.Validator
	store         *mysqlstore.MySQLStore
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	dsn := flag.String("dsn", "p_system_admin:p_system_admin@/p_system?parseTime=true&loc=Local", "MySQL data source name")
	storeKey := flag.String("storekey", "secretkey", "MySQL session store key")
	csrfKey := flag.String("csrfkey", "another-secret-key", "CSRF auth key")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	defer db.Close()

	store, err := newStore(db, *storeKey)
	if err != nil {
		errorLog.Fatal(err)
	}

	defer store.Close()

	templateCache, err := newTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}

	app := &application{
		infoLog:       infoLog,
		errorLog:      errorLog,
		medications:   &models.MedicationModel{DB: db},
		patients:      &models.PatientModel{DB: db},
		users:         &models.UserModel{DB: db},
		templateCache: templateCache,
		decoder:       newDecoder(),
		validator:     validator.NewValidator(),
		store:         store,
	}

	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	srv := &http.Server{
		Addr:         *addr,
		ErrorLog:     errorLog,
		Handler:      app.routes(*csrfKey),
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	errorLog.Fatal(srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem"))
}

// opens and tests the db connection pool before returning it
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

// initiates a new form decoder which will ignore the csrf token input
func newDecoder() *schema.Decoder {
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	return decoder
}

// initiates a mysql session store with a default cleanup interval and registered Flash struct type
func newStore(db *sql.DB, key string) (*mysqlstore.MySQLStore, error) {
	store, err := mysqlstore.NewMySQLStoreFromConnection(db, "sessions", "/", 3600, []byte(key))
	if err != nil {
		return nil, err
	}

	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   864000,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}

	gob.Register(&Flash{})

	store.Cleanup(0)
	return store, nil
}

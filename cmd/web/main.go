package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/go-sql-driver/mysql"
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
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()

	os.Setenv("DBUSER", "admin")
	os.Setenv("DBPASS", "admin")

	cfg := mysql.Config{
		User:                 os.Getenv("DBUSER"),
		Passwd:               os.Getenv("DBPASS"),
		Net:                  "tcp",
		Addr:                 "127.0.0.1:3306",
		DBName:               "p_system",
		AllowNativePasswords: true,
	}

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(cfg)
	if err != nil {
		errorLog.Fatal(err)
	}

	defer db.Close()

	templateCache, err := newTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}

	decoder := schema.NewDecoder()

	app := &application{
		infoLog:       infoLog,
		errorLog:      errorLog,
		medications:   &models.MedicationModel{DB: db},
		patients:      &models.PatientModel{DB: db},
		templateCache: templateCache,
		decoder:       decoder,
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	errorLog.Fatal(srv.ListenAndServe())
}

func openDB(cfg mysql.Config) (*sql.DB, error) {
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

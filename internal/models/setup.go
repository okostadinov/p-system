package models

import "database/sql"

func SetupDB(db *sql.DB) error {
	stmt, err := db.Prepare(`CREATE TABLE IF NOT EXISTS medications (
		name TEXT NOT NULL PRIMARY KEY
	)`)
	if err != nil {
		return err
	}

	_, err = stmt.Exec()
	if err != nil {
		return err
	}

	stmt, err = db.Prepare(`CREATE TABLE IF NOT EXISTS patients (
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		ucn TEXT NOT NULL,
		first_name TEXT NOT NULL,
		last_name TEXT NOT NULL,
		phone_number TEXT NOT NULL,
		height TEXT NOT NULL,
		weight TEXT NOT NULL,
		medication TEXT NOT NULL,
		note TEXT NOT NULL,
		approved BOOLEAN NOT NULL DEFAULT 0,
		first_continuation BOOLEAN NOT NULL DEFAULT 0,
		FOREIGN KEY (medication) REFERENCES medications(name)
		)`)
	if err != nil {
		return err
	}

	_, err = stmt.Exec()
	if err != nil {
		return err
	}

	stmt, err = db.Prepare(`CREATE TABLE IF NOT EXISTS sessions (
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		session_data LONGBLOB,
		created_on TIMESTAMP DEFAULT 0,
		modified_on TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		expires_on TIMESTAMP DEFAULT 0
	)`)
	if err != nil {
		return err
	}

	_, err = stmt.Exec()
	if err != nil {
		return err
	}

	return nil
}

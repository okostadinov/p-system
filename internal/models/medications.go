package models

import "database/sql"

type Medication struct {
	ID   int
	Name string
}

type MedicationModel struct {
	DB *sql.DB
}

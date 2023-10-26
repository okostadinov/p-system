package models

import "database/sql"

type Patient struct {
	ID                int
	Name              string
	UCN               string
	PNumber           string
	Height            int
	Width             int
	Note              string
	MedicationID      int
	Approved          bool
	FirstContinuation bool
}

type PatientModel struct {
	DB *sql.DB
}

func (m *PatientModel) Insert(name string, ucn string, number string, height int, width int, note string) (int, error) {
	return 0, nil
}

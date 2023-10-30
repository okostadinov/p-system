package models

import (
	"database/sql"
)

type Patient struct {
	ID                int
	UCN               int
	Name              string
	PhoneNumber       string
	Height            int
	Weight            int
	Note              string
	MedicationID      int
	Approved          bool
	FirstContinuation bool
}

type PatientModel struct {
	DB *sql.DB
}

func (m *PatientModel) Insert(ucn int, name string, number string, height int, weight int, medicationId int, note string) (int, error) {
	stmt := "INSERT INTO patients (ucn, name, phone_number, height, weight, medication_id, note) VALUES (?, ?, ?, ?, ?, ?, ?)"

	result, err := m.DB.Exec(stmt, ucn, name, number, height, weight, medicationId, note)
	if err != nil {
		return 0, nil
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, nil
	}

	return int(id), err
}

func (m *PatientModel) Get(id int) (*Patient, error) {
	var p Patient

	stmt := "SELECT * FROM patients WHERE id = ?"
	err := m.DB.QueryRow(stmt, id).Scan(&p.ID, &p.UCN, &p.Name, &p.PhoneNumber, &p.Height, &p.Weight, &p.MedicationID, &p.Note, &p.Approved, &p.FirstContinuation)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (m *PatientModel) GetByUCN(ucn int) (*Patient, error) {
	var p Patient

	stmt := "SELECT * FROM patients WHERE ucn = ?"
	err := m.DB.QueryRow(stmt, ucn).Scan(&p.ID, &p.UCN, &p.Name, &p.PhoneNumber, &p.Height, &p.Weight, &p.MedicationID, &p.Note, &p.Approved, &p.FirstContinuation)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (m *PatientModel) Latest() ([]*Patient, error) {
	var patients []*Patient

	stmt := "SELECT * FROM patients ORDER BY ID DECS LIMIT 10"
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var p Patient

		err = rows.Scan(&p.ID, &p.UCN, &p.Name, &p.PhoneNumber, &p.Height, &p.Weight, &p.MedicationID, &p.Note, &p.Approved, &p.FirstContinuation)
		if err != nil {
			return nil, err
		}
		patients = append(patients, &p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return patients, nil
}

func (m *PatientModel) GetAll() ([]*Patient, error) {
	var patients []*Patient

	stmt := "SELECT * FROM patients"
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var p Patient

		err = rows.Scan(&p.ID, &p.UCN, &p.Name, &p.PhoneNumber, &p.Height, &p.Weight, &p.MedicationID, &p.Note, &p.Approved, &p.FirstContinuation)
		if err != nil {
			return nil, err
		}
		patients = append(patients, &p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return patients, nil
}

func (m *PatientModel) GetAllByMedication(medication string) ([]*Patient, error) {
	var patients []*Patient

	stmt := "SELECT patients.* FROM patients JOIN medications ON patients.medication_id = medications.id WHERE medications.name = ?"
	rows, err := m.DB.Query(stmt, medication)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var p Patient

		err = rows.Scan(&p.ID, &p.UCN, &p.Name, &p.PhoneNumber, &p.Height, &p.Weight, &p.MedicationID, &p.Note, &p.Approved, &p.FirstContinuation)
		if err != nil {
			return nil, err
		}
		patients = append(patients, &p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return patients, nil
}

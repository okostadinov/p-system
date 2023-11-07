package models

import (
	"database/sql"
)

type Patient struct {
	ID                int
	UCN               string
	Name              string
	PhoneNumber       string
	Height            int
	Weight            int
	Medication        string
	Note              string
	Approved          bool
	FirstContinuation bool
}

type PatientModel struct {
	DB *sql.DB
}

func (m *PatientModel) Insert(ucn string, name string, number string, height int, weight int, medication string, note string) (int, error) {
	stmt := "INSERT INTO patients (ucn, name, phone_number, height, weight, medication, note) VALUES (?, ?, ?, ?, ?, ?, ?)"

	result, err := m.DB.Exec(stmt, ucn, name, number, height, weight, medication, note)
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
	err := m.DB.QueryRow(stmt, id).Scan(&p.ID, &p.UCN, &p.Name, &p.PhoneNumber, &p.Height, &p.Weight, &p.Medication, &p.Note, &p.Approved, &p.FirstContinuation)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (m *PatientModel) GetByUCN(ucn string) (*Patient, error) {
	var p Patient

	stmt := "SELECT * FROM patients WHERE ucn = ?"
	err := m.DB.QueryRow(stmt, ucn).Scan(&p.ID, &p.UCN, &p.Name, &p.PhoneNumber, &p.Height, &p.Weight, &p.Medication, &p.Note, &p.Approved, &p.FirstContinuation)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (m *PatientModel) Latest() ([]*Patient, error) {
	var patients []*Patient

	stmt := "SELECT * FROM patients ORDER BY ID DESC LIMIT 10"
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var p Patient

		err := rows.Scan(&p.ID, &p.UCN, &p.Name, &p.PhoneNumber, &p.Height, &p.Weight, &p.Medication, &p.Note, &p.Approved, &p.FirstContinuation)
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

		err := rows.Scan(&p.ID, &p.UCN, &p.Name, &p.PhoneNumber, &p.Height, &p.Weight, &p.Medication, &p.Note, &p.Approved, &p.FirstContinuation)
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

	stmt := "SELECT * FROM patients WHERE medication = ?"
	rows, err := m.DB.Query(stmt, medication)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var p Patient

		err := rows.Scan(&p.ID, &p.UCN, &p.Name, &p.PhoneNumber, &p.Height, &p.Weight, &p.Medication, &p.Note, &p.Approved, &p.FirstContinuation)
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

func (m *PatientModel) Update(id int, ucn string, name string, phone string, height int, weight int, medication string, note string, approved bool, firstCont bool) error {
	stmt := "UPDATE patients SET ucn = ?, name = ?, phone_number = ?, height = ?, weight = ?, medication = ?, note = ?, approved = ?, first_continuation = ? WHERE id = ?"

	_, err := m.DB.Exec(stmt, ucn, name, phone, height, weight, medication, note, approved, firstCont, id)
	return err
}

func (m *PatientModel) Delete(id int) error {
	stmt := "DELETE FROM patients WHERE id = ?"

	_, err := m.DB.Exec(stmt, id)
	return err
}

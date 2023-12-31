package models

import (
	"database/sql"
	"errors"
)

type Patient struct {
	ID                int
	UCN               string
	FirstName         string
	LastName          string
	PhoneNumber       string
	Height            int
	Weight            int
	Medication        string
	Note              string
	Approved          bool
	FirstContinuation bool
	UserId            int
}

type PatientModel struct {
	DB *sql.DB
}

func (m *PatientModel) Insert(ucn string, firstName string, lastName string, phone string, height int, weight int, medication string, note string, userId int) (int, error) {
	stmt := "INSERT INTO patients (ucn, first_name, last_name, phone_number, height, weight, medication, note, user_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)"

	result, err := m.DB.Exec(stmt, ucn, firstName, lastName, phone, height, weight, medication, note, userId)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (m *PatientModel) Get(id int) (*Patient, error) {
	var p Patient

	stmt := "SELECT * FROM patients WHERE id = ?"
	err := m.DB.QueryRow(stmt, id).Scan(&p.ID, &p.UCN, &p.FirstName, &p.LastName, &p.PhoneNumber, &p.Height, &p.Weight, &p.Medication, &p.Note, &p.Approved, &p.FirstContinuation, &p.UserId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	return &p, nil
}

func (m *PatientModel) GetByUCN(ucn string) (*Patient, error) {
	var p Patient

	stmt := "SELECT * FROM patients WHERE ucn = ?"
	err := m.DB.QueryRow(stmt, ucn).Scan(&p.ID, &p.UCN, &p.FirstName, &p.LastName, &p.PhoneNumber, &p.Height, &p.Weight, &p.Medication, &p.Note, &p.Approved, &p.FirstContinuation, &p.UserId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
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

		err := rows.Scan(&p.ID, &p.UCN, &p.FirstName, &p.LastName, &p.PhoneNumber, &p.Height, &p.Weight, &p.Medication, &p.Note, &p.Approved, &p.FirstContinuation, &p.UserId)
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

		err := rows.Scan(&p.ID, &p.UCN, &p.FirstName, &p.LastName, &p.PhoneNumber, &p.Height, &p.Weight, &p.Medication, &p.Note, &p.Approved, &p.FirstContinuation, &p.UserId)
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

		err := rows.Scan(&p.ID, &p.UCN, &p.FirstName, &p.LastName, &p.PhoneNumber, &p.Height, &p.Weight, &p.Medication, &p.Note, &p.Approved, &p.FirstContinuation, &p.UserId)
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

func (m *PatientModel) GetAllByUserId(userId int) ([]*Patient, error) {
	var patients []*Patient

	stmt := "SELECT * FROM patients WHERE user_id = ?"
	rows, err := m.DB.Query(stmt, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var p Patient

		err := rows.Scan(&p.ID, &p.UCN, &p.FirstName, &p.LastName, &p.PhoneNumber, &p.Height, &p.Weight, &p.Medication, &p.Note, &p.Approved, &p.FirstContinuation, &p.UserId)
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

func (m *PatientModel) Update(id int, ucn string, firstName string, lastName string, phone string, height int, weight int, medication string, note string, approved bool, firstCont bool, userId int) error {
	stmt := "UPDATE patients SET ucn = ?, first_name = ?, last_name = ?, phone_number = ?, height = ?, weight = ?, medication = ?, note = ?, approved = ?, first_continuation = ? WHERE id = ? && user_id = ?"

	res, err := m.DB.Exec(stmt, ucn, firstName, lastName, phone, height, weight, medication, note, approved, firstCont, id, userId)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrUnauthorizedAction
	}

	return nil
}

func (m *PatientModel) Delete(id int, userId int) error {
	stmt := "DELETE FROM patients WHERE id = ? && user_id = ?"

	res, err := m.DB.Exec(stmt, id, userId)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrUnauthorizedAction
	}

	return nil
}

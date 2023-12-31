package models

import (
	"database/sql"
)

type Medication struct {
	Name   string
	UserId int
}

type MedicationModel struct {
	DB *sql.DB
}

func (m *MedicationModel) Insert(name string, userId int) error {
	stmt := "INSERT INTO medications (name, user_id) VALUES (?, ?)"
	_, err := m.DB.Exec(stmt, name, userId)
	if err != nil {
		return err
	}

	return nil
}

func (m *MedicationModel) GetAll() ([]*Medication, error) {
	var medications []*Medication

	stmt := "SELECT * FROM medications"
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var med Medication

		err = rows.Scan(&med.Name, &med.UserId)
		if err != nil {
			return nil, err
		}
		medications = append(medications, &med)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return medications, nil
}

func (m *MedicationModel) Delete(name string, userId int) error {
	var exists bool
	stmt := "SELECT EXISTS(SELECT true FROM patients WHERE medication = ?)"

	err := m.DB.QueryRow(stmt, name).Scan(&exists)
	if err != nil {
		return err
	}

	if exists {
		return ErrExistingDependency
	}

	stmt = "DELETE FROM medications WHERE name = ? && user_id = ?"

	res, err := m.DB.Exec(stmt, name, userId)
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

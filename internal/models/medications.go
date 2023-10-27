package models

import "database/sql"

type Medication struct {
	ID   int
	Name string
}

type MedicationModel struct {
	DB *sql.DB
}

func (m *MedicationModel) Insert(name string) (int, error) {
	stmt := "INSERT INTO medications (name) VALUES (?)"
	result, err := m.DB.Exec(stmt, name)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (m *MedicationModel) Get(id int) (*Medication, error) {
	var med Medication

	stmt := "SELECT * FROM medications WHERE id = ?"
	err := m.DB.QueryRow(stmt, id).Scan(&med.ID, &med.Name)
	if err != nil {
		return nil, err
	}

	return &med, nil
}

func (m *MedicationModel) GetByName(name string) (*Medication, error) {
	var med Medication

	stmt := "SELECT * FROM medications WHERE name = ?"
	err := m.DB.QueryRow(stmt, name).Scan(&med.ID, &med.Name)
	if err != nil {
		return nil, err
	}

	return &med, nil
}

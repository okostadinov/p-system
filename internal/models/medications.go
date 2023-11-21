package models

import "database/sql"

type Medication struct {
	Name string
}

type MedicationModel struct {
	DB *sql.DB
}

func (m *MedicationModel) Insert(name string) error {
	stmt := "INSERT INTO medications (name) VALUES (?)"
	_, err := m.DB.Exec(stmt, name)
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

		err = rows.Scan(&med.Name)
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

func (m *MedicationModel) Delete(name string) error {
	stmt := "DELETE FROM medications WHERE name = ?"

	_, err := m.DB.Exec(stmt, name)
	return err
}

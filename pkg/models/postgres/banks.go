package postgres

import "database/sql"

type BankModel struct {
	DB *sql.DB
}

func (m *BankModel) Insert(id, name string) error {
	stmt := `
		INSERT INTO bank.banks (id, name)
		VALUES ($1, $2);`
	_, err := m.DB.Exec(stmt, id, name)
	if err != nil {
		return err
	}
	return nil
}

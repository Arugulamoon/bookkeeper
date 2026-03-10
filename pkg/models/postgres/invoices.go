package postgres

import "database/sql"

type InvoicesModel struct {
	DB *sql.DB
}

func (m *InvoicesModel) Insert(
	dueDate, description string, amount int,
) (string, error) {
	stmt := `
		INSERT INTO book.invoice (due_date, description, amount)
		VALUES ($1, $2, $3)
		RETURNING id;`
	var id string
	err := m.DB.QueryRow(stmt, dueDate, description, amount).Scan(&id)
	if err != nil {
		return "", err
	}
	return id, nil
}

package postgres

import (
	"database/sql"
	"time"
)

type JournalEntryModel struct {
	DB *sql.DB
}

func (m *JournalEntryModel) Insert(
	date time.Time, desc string,
) (string, error) {
	stmt := `
		INSERT INTO book.journal_entries (date, description)
		VALUES ($1, $2)
		RETURNING id;`
	var id string
	err := m.DB.QueryRow(stmt, date, desc).Scan(&id)
	if err != nil {
		return "", err
	}
	return id, nil
}

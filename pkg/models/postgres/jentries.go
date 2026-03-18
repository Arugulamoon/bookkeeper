package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type JournalEntryModel struct {
	DB *pgxpool.Pool
}

func (m *JournalEntryModel) Insert(
	ctx context.Context,
	date time.Time, desc string,
) (string, error) {
	stmt := `
		INSERT INTO book.journal_entries (date, description)
		VALUES ($1, $2)
		RETURNING id;`
	var id string
	err := m.DB.QueryRow(ctx, stmt, date, desc).Scan(&id)
	if err != nil {
		return "", err
	}
	return id, nil
}

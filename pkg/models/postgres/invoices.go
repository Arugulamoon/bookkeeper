package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type InvoicesModel struct {
	DB *pgxpool.Pool
}

func (m *InvoicesModel) Insert(
	ctx context.Context,
	dueDate, description string, amount int,
) (string, error) {
	stmt := `
		INSERT INTO book.invoice (due_date, description, amount)
		VALUES ($1, $2, $3)
		RETURNING id;`
	var id string
	err := m.DB.QueryRow(ctx, stmt, dueDate, description, amount).Scan(&id)
	if err != nil {
		return "", err
	}
	return id, nil
}

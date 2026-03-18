package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type BankModel struct {
	DB *pgxpool.Pool
}

func (m *BankModel) Insert(
	ctx context.Context,
	id, name string,
) error {
	stmt := `
		INSERT INTO bank.banks (id, name)
		VALUES ($1, $2);`
	_, err := m.DB.Exec(ctx, stmt, id, name)
	if err != nil {
		return err
	}
	return nil
}

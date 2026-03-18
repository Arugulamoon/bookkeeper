package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Arugulamoon/bookkeeper/pkg/models"
)

type BankCurrencyModel struct {
	DB *pgxpool.Pool
}

func (m *BankCurrencyModel) SelectAll(
	ctx context.Context,
) ([]models.Currency, error) {
	stmt := `SELECT id, name FROM bank.currencies;`
	rows, err := m.DB.Query(ctx, stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var currencies []models.Currency
	for rows.Next() {
		var currency models.Currency
		err := rows.Scan(&currency.Id, &currency.Name)
		if err != nil {
			return nil, err
		}
		currencies = append(currencies, currency)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return currencies, nil
}

func (m *BankCurrencyModel) Insert(
	ctx context.Context,
	id, name string,
) (int, error) {
	stmt := `INSERT INTO bank.currencies (id, name) VALUES ($1, $2);`
	res, err := m.DB.Exec(ctx, stmt, id, name)
	if err != nil {
		return 0, err
	}
	return int(res.RowsAffected()), nil
}

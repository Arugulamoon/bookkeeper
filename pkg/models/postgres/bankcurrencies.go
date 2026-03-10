package postgres

import (
	"database/sql"

	"github.com/Arugulamoon/bookkeeper/pkg/models"
)

type BankCurrencyModel struct {
	DB *sql.DB
}

func (m *BankCurrencyModel) SelectAll() ([]models.Currency, error) {
	stmt := `SELECT id, name FROM bank.currencies;`
	rows, err := m.DB.Query(stmt)
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

func (m *BankCurrencyModel) Insert(id, name string) (int, error) {
	stmt := `INSERT INTO bank.currencies (id, name) VALUES ($1, $2);`
	res, err := m.DB.Exec(stmt, id, name)
	if err != nil {
		return 0, err
	}

	numInserted, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return int(numInserted), nil
}

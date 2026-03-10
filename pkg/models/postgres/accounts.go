package postgres

import (
	"database/sql"
	"errors"

	"github.com/Arugulamoon/bookkeeper/pkg/models"
)

type AccountModel struct {
	DB *sql.DB
}

func (m *AccountModel) SelectAll() ([]*models.Account, error) {
	stmt := `
		SELECT
			accts.account_type,
			accts.name,
			accts.bank_account_id,
			accts.sort_order
		FROM book.accounts AS accts
		INNER JOIN book.account_types AS accttypes
			ON accts.account_type = accttypes.name
		ORDER BY accttypes.sort_order, accttypes.name, accts.sort_order, accts.name;`
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	accts := []*models.Account{}
	for rows.Next() {
		acct := &models.Account{}
		err := rows.Scan(
			&acct.AccountType,
			&acct.Name,
			&acct.BankAccountId,
			&acct.SortOrder,
		)
		if err != nil {
			return nil, err
		}
		accts = append(accts, acct)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return accts, nil
}

func (m *AccountModel) SelectAllNamesByAccountType(
	acctType string,
) ([]string, error) {
	stmt := `
		SELECT name
		FROM book.accounts
		WHERE account_type = $1
		ORDER BY sort_order, name;`
	rows, err := m.DB.Query(stmt, acctType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	names := []string{}
	for rows.Next() {
		name := ""
		err := rows.Scan(&name)
		if err != nil {
			return nil, err
		}
		names = append(names, name)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return names, nil
}

func (m *AccountModel) Select(
	acctType, acctName string,
) (*models.Account, error) {
	stmt := `
		SELECT account_type, name, bank_account_id, sort_order
		FROM book.accounts
		WHERE account_type = $1 AND name = $2;`
	acct := &models.Account{}
	err := m.DB.QueryRow(stmt, acctType, acctName).Scan(
		&acct.AccountType,
		&acct.Name,
		&acct.BankAccountId,
		&acct.SortOrder,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}
	return acct, nil
}

func (m *AccountModel) SelectByBankAccountId(
	bankAcctId string,
) (*models.Account, error) {
	stmt := `
		SELECT account_type, name, bank_account_id, sort_order
		FROM book.accounts
		WHERE bank_account_id = $1;`
	acct := &models.Account{}
	err := m.DB.QueryRow(stmt, bankAcctId).Scan(
		&acct.AccountType,
		&acct.Name,
		&acct.BankAccountId,
		&acct.SortOrder,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}
	return acct, nil
}

func (m *AccountModel) Insert(
	acctType, name string, bankAccountId *string, sortOrder int,
) (int, error) {
	stmt := `
		INSERT INTO book.accounts
			(account_type, name, bank_account_id, sort_order)
		VALUES ($1, $2, $3, $4);`
	res, err := m.DB.Exec(stmt, acctType, name, bankAccountId, sortOrder)
	if err != nil {
		return 0, err
	}

	numInserted, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return int(numInserted), nil
}

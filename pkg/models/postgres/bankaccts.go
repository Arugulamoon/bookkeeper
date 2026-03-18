package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Arugulamoon/bookkeeper/pkg/models"
)

type BankAccountModel struct {
	DB *pgxpool.Pool
}

func (m *BankAccountModel) GetId(
	ctx context.Context,
	name string,
) (*string, error) {
	stmt := `SELECT id FROM bank.accounts WHERE name = $1;`
	var id *string
	err := m.DB.QueryRow(ctx, stmt, name).Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	return id, nil
}

func (m *BankAccountModel) Insert(
	ctx context.Context,
	name, bankId, acctType string,
) (string, error) {
	stmt := `
		INSERT INTO bank.accounts (name, bank_id, account_type)
		VALUES ($1, $2, $3)
		RETURNING id;`
	var id string
	err := m.DB.QueryRow(ctx, stmt, name, bankId, acctType).Scan(&id)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (m *BankAccountModel) SelectPaymentDescriptionsByPaymentType(
	ctx context.Context,
	acctId, paymentType string,
) ([]string, error) {
	stmt := `
		SELECT description
		FROM bank.account_payment_descriptions
		WHERE account_id = $1 AND payment_type = $2;`
	rows, err := m.DB.Query(ctx, stmt, acctId, paymentType)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var descs []string
	for rows.Next() {
		var desc string
		err := rows.Scan(&desc)
		if err != nil {
			return nil, err
		}
		descs = append(descs, desc)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return descs, nil
}

func (m *BankAccountModel) InsertPaymentDescription(
	ctx context.Context,
	acctId, paymentType, desc string,
) error {
	stmt := `
		INSERT INTO bank.account_payment_descriptions
			(account_id, payment_type, description)
		VALUES ($1, $2, $3);`
	_, err := m.DB.Exec(ctx, stmt, acctId, paymentType, desc)
	if err != nil {
		return err
	}
	return nil
}

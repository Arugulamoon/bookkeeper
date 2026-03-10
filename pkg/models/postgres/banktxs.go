package postgres

import (
	"database/sql"
	"time"

	"github.com/Arugulamoon/bookkeeper/pkg/models"
)

const SELECT_BANK_TRANSACTIONS_STMT_PREFIX = `
	SELECT
		txs.id,
		txs.date,
		txs.description,
		txs.debit,
		txs.credit,
		accts.id AS account_id
	FROM bank.transactions AS txs
	INNER JOIN bank.accounts AS accts
		ON txs.account_id = accts.id`

type BankTransactionModel struct {
	DB *sql.DB
}

func (m *BankTransactionModel) SelectAllCreditCardPaymentsReceived() ([]*models.BankTransaction, error) {
	stmt := SELECT_BANK_TRANSACTIONS_STMT_PREFIX + `
		WHERE
			txs.currency_id = 'CAD' AND
			txs.description ILIKE ANY (ARRAY['PRE-AUTHORIZED PAYMENT%', 'AUTOMATIC PAYMENT%', 'PAYMENT - THANK YOU%'])
		ORDER BY txs.date;`
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanBankTransactionRows(rows)
}

func (m *BankTransactionModel) SelectAllPaymentsMadeToCreditCard() ([]*models.BankTransaction, error) {
	stmt := SELECT_BANK_TRANSACTIONS_STMT_PREFIX + `
		WHERE
			txs.currency_id = 'CAD' AND
			txs.description ILIKE ANY (ARRAY['MISC PAYMENT RBC CREDIT CARD', 'MISC PAYMENT CIBC CPD'])
		ORDER BY txs.date;`
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanBankTransactionRows(rows)
}

func (m *BankTransactionModel) SelectAllPaymentsMadeToOpaqueCreditCard() ([]*models.BankTransaction, error) {
	stmt := SELECT_BANK_TRANSACTIONS_STMT_PREFIX + `
		LEFT OUTER JOIN book.journal_entry_account_entries AS jentry
			ON txs.id = jentry.bank_transaction_id
		WHERE
			txs.currency_id = 'CAD' AND
			txs.description ILIKE ANY (ARRAY['MISC PAYMENT RBC CREDIT CARD', 'MISC PAYMENT CIBC CPD']) AND
			jentry.bank_transaction_id IS NULL
		ORDER BY txs.date;`
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanBankTransactionRows(rows)
}

func (m *BankTransactionModel) SelectAllNonCreditCardPayments() ([]*models.BankTransaction, error) {
	stmt := SELECT_BANK_TRANSACTIONS_STMT_PREFIX + `
		WHERE
			txs.currency_id = 'CAD' AND
			txs.description NOT ILIKE ALL (ARRAY['MISC PAYMENT RBC CREDIT CARD', 'MISC PAYMENT CIBC CPD', 'PRE-AUTHORIZED PAYMENT%', 'AUTOMATIC PAYMENT%', 'PAYMENT - THANK YOU%'])
		ORDER BY txs.date;`
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanBankTransactionRows(rows)
}

func scanBankTransactionRows(rows *sql.Rows) ([]*models.BankTransaction, error) {
	txs := []*models.BankTransaction{}
	for rows.Next() {
		tx := &models.BankTransaction{}
		err := rows.Scan(
			&tx.Id,
			&tx.Date,
			&tx.Description,
			&tx.Debit,
			&tx.Credit,
			&tx.AccountId,
		)
		if err != nil {
			return nil, err
		}
		txs = append(txs, tx)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return txs, nil
}

func (m *BankTransactionModel) InsertRBC(
	date time.Time,
	desc, desc2 string,
	debit, credit float64,
	currency string,
	acctNum, chequeNum string,
	acctId string,
) (string, error) {
	stmt := `
		INSERT INTO bank.transactions
			(date, description, description2, debit, credit, currency_id, account_number, cheque_number, account_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id;`

	var id string
	err := m.DB.QueryRow(stmt,
		date.Format("2006-01-02"),
		desc,
		desc2,
		debit,
		credit,
		currency,
		acctNum,
		chequeNum,
		acctId,
	).Scan(&id)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (m *BankTransactionModel) InsertCIBC(
	date time.Time,
	desc string,
	debit, credit float64,
	cardNum string,
	acctId string,
) (string, error) {
	var id string
	stmt := `
		INSERT INTO bank.transactions
			(date, description, debit, credit, currency_id, card_number, account_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id;`
	err := m.DB.QueryRow(stmt,
		date.Format("2006-01-02"),
		desc,
		debit,
		credit,
		"CAD",
		cardNum,
		acctId,
	).Scan(&id)
	if err != nil {
		return "", err
	}
	return id, nil
}

package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/Arugulamoon/bookkeeper/pkg/models"
)

const SELECT_STATEMENT_PREFIX = `
	SELECT
		acctentry.id,
		jentry.id AS journal_entry_id,
		jentry.date,
		jentry.description,
		acctentry.balance_type,
		acctentry.assigner_id,
		acctentry.account_type,
		acctentry.account_name,
		accttypes.balance_type AS account_balance_type,
		acctentry.amount,
		bankaccts.id AS bank_account_id,
		bankaccts.name AS bank_account_name,
		acctentry.bank_transaction_id
	FROM book.journal_entry_account_entries AS acctentry
	INNER JOIN book.journal_entries AS jentry
		ON acctentry.journal_entry_id = jentry.id
	INNER JOIN book.accounts AS accts
		ON
			acctentry.account_type = accts.account_type AND
			acctentry.account_name = accts.name
	INNER JOIN book.account_types AS accttypes
		ON accts.account_type = accttypes.name
	INNER JOIN bank.transactions AS txs
		ON acctentry.bank_transaction_id = txs.id
	INNER JOIN bank.accounts AS bankaccts
		ON txs.account_id = bankaccts.id`

type JournalAccountEntryModel struct {
	DB *sql.DB
}

func (m *JournalAccountEntryModel) SelectAllByAssignerId(
	assignerId string,
) ([]*models.JournalAccountEntry, error) {
	stmt := SELECT_STATEMENT_PREFIX + `
		WHERE acctentry.assigner_id = $1
		ORDER BY jentry.date DESC;`
	rows, err := m.DB.Query(stmt, assignerId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanRows(rows)
}

func (m *JournalAccountEntryModel) SelectAllByAccountId(
	acctType, acctName string,
) ([]*models.JournalAccountEntry, error) {
	stmt := SELECT_STATEMENT_PREFIX + `
		WHERE acctentry.account_type = $1 AND acctentry.account_name = $2
		ORDER BY jentry.date DESC;`
	rows, err := m.DB.Query(stmt, acctType, acctName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanRows(rows)
}

func (m *JournalAccountEntryModel) SelectAllByLikeDescription(
	str string,
) ([]*models.JournalAccountEntry, error) {
	stmt := SELECT_STATEMENT_PREFIX + `
		WHERE
			accts.account_type = 'Expense' AND
			jentry.description ILIKE CONCAT('%', $1::text, '%');`
	rows, err := m.DB.Query(stmt, str)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanRows(rows)
}

func (m *JournalAccountEntryModel) SelectAllByAccountForYear(
	acctType, acctName string,
	year int,
) ([]*models.JournalAccountEntry, error) {
	startDate, err := time.Parse("2006-01-02", fmt.Sprintf("%d-01-01", year))
	if err != nil {
		return nil, err
	}

	endDate := startDate.AddDate(1, 0, -1) // add a year and subtract a day

	return m.SelectAllByAccountForDateRange(acctType, acctName,
		startDate.Format("2006-01-02"),
		endDate.Format("2006-01-02"),
	)
}

func (m *JournalAccountEntryModel) SelectAllByAccountForMonth(
	acctType, acctName string,
	year, month int,
) ([]*models.JournalAccountEntry, error) {
	var monthStr string
	if month < 10 {
		monthStr = fmt.Sprintf("0%d", month)
	} else {
		monthStr = strconv.Itoa(month)
	}
	startDate, err := time.Parse("2006-01-02", fmt.Sprintf("%d-%s-01", year, monthStr))
	if err != nil {
		return nil, err
	}

	endDate := startDate.AddDate(0, 1, -1) // add a month and subtract a day

	return m.SelectAllByAccountForDateRange(acctType, acctName,
		startDate.Format("2006-01-02"),
		endDate.Format("2006-01-02"),
	)
}

func (m *JournalAccountEntryModel) SelectAllByAccountForDateRange(
	acctType, acctName string,
	startDate, endDate string,
) ([]*models.JournalAccountEntry, error) {
	stmt := SELECT_STATEMENT_PREFIX + `
		WHERE
			acctentry.account_type = $1 AND
			acctentry.account_name = $2 AND
			jentry.date >= $3::date AND
			jentry.date <= $4::date
		ORDER BY jentry.date, jentry.description;`
	rows, err := m.DB.Query(stmt, acctType, acctName, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanRows(rows)
}

func scanRows(rows *sql.Rows) ([]*models.JournalAccountEntry, error) {
	entries := []*models.JournalAccountEntry{}
	for rows.Next() {
		entry := &models.JournalAccountEntry{}
		err := rows.Scan(
			&entry.Id,
			&entry.JournalEntryId,
			&entry.Date,
			&entry.Description,
			&entry.BalanceType,
			&entry.AssignerId,
			&entry.AccountType,
			&entry.AccountName,
			&entry.AccountBalanceType,
			&entry.Amount,
			&entry.BankAccountId,
			&entry.BankAccountName,
			&entry.BankTransactionId,
		)
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return entries, nil
}

func (m *JournalAccountEntryModel) SelectAllAlreadyImportedBankTransactionIds() ([]string, error) {
	stmt := `
		SELECT DISTINCT bank_transaction_id
		FROM book.journal_entry_account_entries;`
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	ids := make([]string, 0)
	for rows.Next() {
		var id string
		err := rows.Scan(&id)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return ids, nil
}

func (m *JournalAccountEntryModel) Select(id string) (*models.JournalAccountEntry, error) {
	stmt := SELECT_STATEMENT_PREFIX + `
		WHERE acctentry.id = $1;`
	entry := &models.JournalAccountEntry{}
	err := m.DB.QueryRow(stmt, id).Scan(
		&entry.Id,
		&entry.JournalEntryId,
		&entry.Date,
		&entry.Description,
		&entry.BalanceType,
		&entry.AssignerId,
		&entry.AccountType,
		&entry.AccountName,
		&entry.AccountBalanceType,
		&entry.Amount,
		&entry.BankAccountId,
		&entry.BankAccountName,
		&entry.BankTransactionId,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	return entry, nil
}

func (m *JournalAccountEntryModel) Insert(
	jEntryId, balType string,
	assignerId *string,
	acctType, acctName string,
	amount float64,
	bankTxId string,
) (string, error) {
	stmt := `
		INSERT INTO book.journal_entry_account_entries
			(journal_entry_id, balance_type, assigner_id, account_type, account_name, amount, bank_transaction_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id;`

	var id string
	err := m.DB.QueryRow(stmt, jEntryId, balType, assignerId, acctType, acctName, amount, bankTxId).Scan(&id)
	if err != nil {
		return "", err
	}

	return id, nil
}

// TODO (10): Implement Batch
func (m *JournalAccountEntryModel) UpdateAccount(id, acctType, acctName string) (string, error) {
	stmt := `
		UPDATE book.journal_entry_account_entries
		SET account_type = $2, account_name = $3
		WHERE id = $1;`
	res, err := m.DB.Exec(stmt, id, acctType, acctName)
	if err != nil {
		return "", err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return "", err
	}

	if rowsAffected != 1 {
		return "", fmt.Errorf("incorrect number of rows (%d) affected", rowsAffected)
	}

	return id, nil
}

// TODO (10): Implement Batch
func (m *JournalAccountEntryModel) UpdateAmount(id string, amount float64) (string, error) {
	stmt := `
		UPDATE book.journal_entry_account_entries
		SET amount = $2
		WHERE id = $1;`
	res, err := m.DB.Exec(stmt, id, amount)
	if err != nil {
		return "", err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return "", err
	}

	if rowsAffected != 1 {
		return "", fmt.Errorf("incorrect number of rows (%d) affected", rowsAffected)
	}

	return id, nil
}

package postgres

import (
	"database/sql"

	"github.com/Arugulamoon/bookkeeper/pkg/models"
)

type AssignerModel struct {
	DB *sql.DB
}

type assignerDesc struct {
	Id                         string
	Name                       string
	BankTransactionDescription string
	AccountType                string
	AccountName                string
}

func (m *AssignerModel) SelectAll() ([]*models.Assigner, error) {
	stmt := `
		SELECT
			assigners.id,
			assigners.name,
			assigners.account_type,
			assigners.account_name,
			descs.bank_transaction_description
		FROM book.assigner_bank_transaction_descriptions AS descs
		INNER JOIN book.assigners AS assigners
			ON descs.assigner_id = assigners.id
		ORDER BY assigners.name;`
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	assignerDescs := []*assignerDesc{}
	for rows.Next() {
		assignerDesc := &assignerDesc{}
		err := rows.Scan(
			&assignerDesc.Id,
			&assignerDesc.Name,
			&assignerDesc.AccountType,
			&assignerDesc.AccountName,
			&assignerDesc.BankTransactionDescription,
		)
		if err != nil {
			return nil, err
		}
		assignerDescs = append(assignerDescs, assignerDesc)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	// flatten descriptions into slice of descs
	assigners := []*models.Assigner{}
	for _, assignerDesc := range assignerDescs {
		foundIdx := -1
		for i, assigner := range assigners {
			if assigner.Id == assignerDesc.Id {
				foundIdx = i
				break
			}
		}
		if foundIdx > -1 {
			assigners[foundIdx].BankTransactionDescriptions =
				append(assigners[foundIdx].BankTransactionDescriptions,
					assignerDesc.BankTransactionDescription)
		} else {
			assigners = append(assigners, &models.Assigner{
				Id:                          assignerDesc.Id,
				Name:                        assignerDesc.Name,
				AccountType:                 assignerDesc.AccountType,
				AccountName:                 assignerDesc.AccountName,
				BankTransactionDescriptions: []string{assignerDesc.BankTransactionDescription},
			})
		}
	}

	return assigners, nil
}

func (m *AssignerModel) Select(id string) (*models.Assigner, error) {
	stmt := `
		SELECT
			assigners.id,
			assigners.name,
			assigners.account_type,
			assigners.account_name,
			descs.bank_transaction_description
		FROM book.assigner_bank_transaction_descriptions AS descs
		INNER JOIN book.assigners AS assigners
			ON descs.assigner_id = assigners.id
		WHERE assigners.id = $1;`
	rows, err := m.DB.Query(stmt, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	assignerDescs := []*assignerDesc{}
	for rows.Next() {
		assignerDesc := &assignerDesc{}
		err := rows.Scan(
			&assignerDesc.Id,
			&assignerDesc.Name,
			&assignerDesc.AccountType,
			&assignerDesc.AccountName,
			&assignerDesc.BankTransactionDescription,
		)
		if err != nil {
			return nil, err
		}
		assignerDescs = append(assignerDescs, assignerDesc)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	// flatten descriptions into slice of descs
	var assigner *models.Assigner
	for i, assignerDesc := range assignerDescs {
		if i == 0 {
			assigner = &models.Assigner{
				Id:                          assignerDesc.Id,
				Name:                        assignerDesc.Name,
				AccountType:                 assignerDesc.AccountType,
				AccountName:                 assignerDesc.AccountName,
				BankTransactionDescriptions: []string{assignerDesc.BankTransactionDescription},
			}
		} else {
			assigner.BankTransactionDescriptions = append(assigner.BankTransactionDescriptions,
				assignerDesc.BankTransactionDescription)
		}
	}

	return assigner, nil
}

func (m *AssignerModel) Insert(name, acctType, acctName string) (string, error) {
	stmt := `
		INSERT INTO book.assigners (name, account_type, account_name)
		VALUES ($1, $2, $3)
		RETURNING id;`
	var id string
	err := m.DB.QueryRow(stmt, name, acctType, acctName).Scan(&id)
	if err != nil {
		return "", err
	}

	return id, nil
}

// TODO: Replace with Insert
func (m *AssignerModel) InsertBankTransactionDescription(
	desc, assignerId string,
) (string, error) {
	stmt := `
		INSERT INTO book.assigner_bank_transaction_descriptions
			(bank_transaction_description, assigner_id)
		VALUES ($1, $2)
		RETURNING id;`
	var id string
	err := m.DB.QueryRow(stmt, desc, assignerId).Scan(&id)
	if err != nil {
		return "", err
	}

	return id, nil
}

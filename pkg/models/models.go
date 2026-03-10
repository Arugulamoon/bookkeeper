package models

import (
	"database/sql"
	"errors"
	"time"
)

const UUID_MAX_LENGTH = 36

var ErrNoRecord = errors.New("models: no matching record found")

type Currency struct {
	Id   string
	Name string
}

type AccountAnnualTotal struct {
	Account       *Account
	MonthlyTotals []float64 // slice index = month - 1
	AnnualTotal   float64
}

type AccountTypeAnnualTotal struct {
	MonthlyTotals []float64 // slice index = month - 1
	AnnualTotal   float64
}

type JournalAccountEntry struct {
	Id                 string
	JournalEntryId     string
	Date               time.Time
	Description        string
	BalanceType        string
	AssignerId         *string
	AccountType        string
	AccountName        string
	AccountBalanceType string
	Amount             float64
	BankAccountId      string
	BankAccountName    string
	BankTransactionId  string
}

type Assigner struct {
	Id                          string
	Name                        string
	AccountType                 string
	AccountName                 string
	BankTransactionDescriptions []string
}

type Account struct {
	AccountType   string
	Name          string
	BankAccountId sql.NullString
	SortOrder     int
}

package models

import "time"

type BankTransaction struct {
	Id          string
	Date        time.Time
	Description string
	Debit       float64
	Credit      float64
	AccountId   string
}

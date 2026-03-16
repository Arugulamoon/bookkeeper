package models

import "time"

type SchoolExpense struct {
	Id              string
	DueDate         time.Time
	Description     string
	Amount          int
	InvoiceId       string
	SchoolYear      string
	SchoolId        string
	GradeId         string
	EventId         string
	DatePaid        time.Time
	EventMarkedPaid bool
}

package postgres

import "database/sql"

type SchoolExpensesModel struct {
	DB *sql.DB
}

func (m *SchoolExpensesModel) InsertInvoice(
	invoiceId string,
	schoolYear, school, grade string,
	eventId string,
	datePaid *string,
	eventMarkedPaid bool,
) (string, error) {
	stmt := `
		INSERT INTO school.invoice
			(invoice_id, school_year, school_id, grade_id, event_id, date_paid, event_marked_paid)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id;`
	var id string
	err := m.DB.QueryRow(stmt,
		invoiceId,
		schoolYear,
		school,
		grade,
		eventId,
		datePaid,
		eventMarkedPaid,
	).Scan(&id)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (m *SchoolExpensesModel) UpdateInvoiceEventId(
	id, eventId string,
) (int, error) {
	stmt := `
		UPDATE school.invoice
		SET event_id = $2
		WHERE id = $1;`
	res, err := m.DB.Exec(stmt, id, eventId)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return int(rowsAffected), nil
}

func (m *SchoolExpensesModel) UpdateInvoiceEventMarkedPaid(
	id string,
) (int, error) {
	stmt := `
		UPDATE school.invoice
		SET event_marked_paid = TRUE
		WHERE id = $1;`
	res, err := m.DB.Exec(stmt, id)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return int(rowsAffected), nil
}

func (m *SchoolExpensesModel) InsertReimbursement(
	invoiceId, split string, amount *int, date *string,
) (string, error) {
	stmt := `
		INSERT INTO school.reimbursement
			(invoice_id, split, amount, date)
		VALUES ($1, $2, $3, $4)
		RETURNING id;`
	var id string
	err := m.DB.QueryRow(stmt, invoiceId, split, amount, date).Scan(&id)
	if err != nil {
		return "", err
	}
	return id, nil
}

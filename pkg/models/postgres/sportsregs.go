package postgres

import (
	"database/sql"

	"github.com/Arugulamoon/bookkeeper/pkg/models"
)

type SportsRegistrationsModel struct {
	DB *sql.DB
}

func (m *SportsRegistrationsModel) SelectAll() ([]*models.SportsRegistration, error) {
	stmt := `
		SELECT *
		FROM sports.registration AS regs
		ORDER BY regs.start_date, regs.end_date;`
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanSportsRegistrationRows(rows)
}

func (m *SportsRegistrationsModel) SelectAllCurrentAndUpcoming() ([]*models.SportsRegistration, error) {
	stmt := `
		SELECT *
		FROM sports.registration AS regs
		WHERE regs.end_date > NOW()
		ORDER BY regs.start_date, regs.end_date;`
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanSportsRegistrationRows(rows)
}

func (m *SportsRegistrationsModel) SelectAllPast() ([]*models.SportsRegistration, error) {
	stmt := `
		SELECT *
		FROM sports.registration AS regs
		WHERE regs.end_date < NOW()
		ORDER BY regs.start_date, regs.end_date;`
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanSportsRegistrationRows(rows)
}

func scanSportsRegistrationRows(rows *sql.Rows) ([]*models.SportsRegistration, error) {
	entries := []*models.SportsRegistration{}
	for rows.Next() {
		entry := &models.SportsRegistration{}
		err := rows.Scan(
			&entry.Id,
			&entry.Name,
			&entry.PriceInCents,
			&entry.RegularPriceInCents,
			&entry.DiscountInCents,
			&entry.TaxInCents,
			&entry.Location,
			&entry.Day,
			&entry.StartTime,
			&entry.StartTimeRange,
			&entry.EndTimeRange,
			&entry.Duration,
			&entry.StartDate,
			&entry.EndDate,
			&entry.Sessions,
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

func (m *SportsRegistrationsModel) Insert(
	name string,
	totalPrice, regularPrice, discount, tax int,
	location string,
	day, startTime, startTimeRange, endTimeRange string, duration int,
	startDate, endDate string,
	sessions int,
) (int, error) {
	stmt := `
		INSERT INTO sports.registration
			(name, price, regular_price, discount, tax, location, day, start_time, start_time_range, end_time_range, duration, start_date, end_date, sessions)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14);`
	res, err := m.DB.Exec(stmt,
		name,
		totalPrice,
		regularPrice,
		discount,
		tax,
		location,
		day,
		startTime,
		startTimeRange,
		endTimeRange,
		duration,
		startDate,
		endDate,
		sessions,
	)
	if err != nil {
		return 0, err
	}

	numInserted, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return int(numInserted), nil
}

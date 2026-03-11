package postgres

import (
	"database/sql"
	"time"
)

type SportsMembershipModel struct {
	DB *sql.DB
}

func (m *SportsMembershipModel) Insert(
	name, seasonYear, seasonType, location string,
) (string, error) {
	stmt := `
		INSERT INTO sports.memberships
			(name, season_year, season_type, location)
		VALUES ($1, $2, $3, $4)
		RETURNING id;`
	var id string
	err := m.DB.QueryRow(stmt, name, seasonYear, seasonType, location).Scan(&id)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (m *SportsMembershipModel) InsertHomeGame(
	membershipId string,
	date time.Time,
	startTime string,
	opponent string,
	notes string,
	location string,
	eventId string,
) (string, error) {
	stmt := `
		INSERT INTO sports.membership_games
			(membership_id, date, start_time, opponent, notes, location, event_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id;`
	var id string
	err := m.DB.QueryRow(stmt,
		membershipId,
		date,
		startTime,
		opponent,
		notes,
		location,
		eventId,
	).Scan(&id)
	if err != nil {
		return "", err
	}
	return id, nil
}

package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Arugulamoon/bookkeeper/pkg/models"
)

type SportsMembershipModel struct {
	DB *pgxpool.Pool
}

func (m *SportsMembershipModel) SelectAll(
	ctx context.Context,
) ([]*models.SportsMembership, error) {
	stmt := `
		SELECT id, name, season_year, season_type, location
		FROM sports.memberships
		ORDER BY season_year, season_type DESC, name;`
	rows, err := m.DB.Query(ctx, stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	memberships := []*models.SportsMembership{}
	for rows.Next() {
		membership := &models.SportsMembership{}
		err := rows.Scan(
			&membership.Id,
			&membership.Name,
			&membership.SeasonYear,
			&membership.SeasonType,
			&membership.Location,
		)
		if err != nil {
			return nil, err
		}
		memberships = append(memberships, membership)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return memberships, nil
}

func (m *SportsMembershipModel) Select(
	ctx context.Context,
	id string,
) (*models.SportsMembership, error) {
	stmt := `
		SELECT
			id,
			name,
			season_year,
			season_type,
			location
		FROM sports.memberships
		WHERE id = $1;`
	membership := &models.SportsMembership{}
	err := m.DB.QueryRow(ctx, stmt, id).Scan(
		&membership.Id,
		&membership.Name,
		&membership.SeasonYear,
		&membership.SeasonType,
		&membership.Location,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}
	return membership, nil
}

func (m *SportsMembershipModel) Insert(
	ctx context.Context,
	name, seasonYear, seasonType, location string,
) (string, error) {
	stmt := `
		INSERT INTO sports.memberships
			(name, season_year, season_type, location)
		VALUES ($1, $2, $3, $4)
		RETURNING id;`
	var id string
	err := m.DB.QueryRow(ctx, stmt, name, seasonYear, seasonType, location).Scan(&id)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (m *SportsMembershipModel) SelectAllGames(
	ctx context.Context,
	id string,
) ([]*models.SportsMembershipGame, error) {
	stmt := `
		SELECT
			id,
			date,
			start_time,
			opponent,
			notes,
			location,
			event_id
		FROM sports.membership_games
		WHERE membership_id = $1
		ORDER BY date;`
	rows, err := m.DB.Query(ctx, stmt, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	games := []*models.SportsMembershipGame{}
	for rows.Next() {
		game := &models.SportsMembershipGame{}
		err := rows.Scan(
			&game.Id,
			&game.Date,
			&game.StartTime,
			&game.Opponent,
			&game.Notes,
			&game.Location,
			&game.EventId,
		)
		if err != nil {
			return nil, err
		}
		games = append(games, game)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return games, nil
}

func (m *SportsMembershipModel) InsertGame(
	ctx context.Context,
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
	err := m.DB.QueryRow(ctx, stmt,
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

func (m *SportsMembershipModel) UpdateGameEventId(
	ctx context.Context,
	id, eventId string,
) (int, error) {
	stmt := `
		UPDATE sports.membership_games
		SET event_id = $2
		WHERE id = $1;`
	res, err := m.DB.Exec(ctx, stmt, id, eventId)
	if err != nil {
		return 0, err
	}
	return int(res.RowsAffected()), nil
}

package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type SchoolModel struct {
	DB *pgxpool.Pool
}

func (m *SchoolModel) InsertGrade(
	ctx context.Context,
	id, name string,
) error {
	stmt := `
		INSERT INTO school.grade (id, name)
		VALUES ($1, $2);`
	_, err := m.DB.Exec(ctx, stmt, id, name)
	if err != nil {
		return err
	}
	return nil
}

func (m *SchoolModel) InsertSchool(
	ctx context.Context,
	id, name string, address, phone, principal *string,
) error {
	stmt := `
		INSERT INTO school.school (id, name, address, phone, principal)
		VALUES ($1, $2, $3, $4, $5);`
	_, err := m.DB.Exec(ctx, stmt, id, name, address, phone, principal)
	if err != nil {
		return err
	}
	return nil
}

func (m *SchoolModel) InsertSchoolYear(
	ctx context.Context,
	schoolYear, schoolId, gradeId string, teacher, education *string,
) error {
	stmt := `
		INSERT INTO school.school_year
			(school_year, school_id, grade_id, teacher, education)
		VALUES ($1, $2, $3, $4, $5);`
	_, err := m.DB.Exec(ctx, stmt,
		schoolYear,
		schoolId,
		gradeId,
		teacher,
		education,
	)
	if err != nil {
		return err
	}
	return nil
}

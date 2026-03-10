package postgres

import "database/sql"

type SchoolModel struct {
	DB *sql.DB
}

func (m *SchoolModel) InsertGrade(id, name string) error {
	stmt := `
		INSERT INTO school.grade (id, name)
		VALUES ($1, $2);`
	_, err := m.DB.Exec(stmt, id, name)
	if err != nil {
		return err
	}
	return nil
}

func (m *SchoolModel) InsertSchool(
	id, name string, address, phone, principal *string,
) error {
	stmt := `
		INSERT INTO school.school (id, name, address, phone, principal)
		VALUES ($1, $2, $3, $4, $5);`
	_, err := m.DB.Exec(stmt, id, name, address, phone, principal)
	if err != nil {
		return err
	}
	return nil
}

func (m *SchoolModel) InsertSchoolYear(
	schoolYear, schoolId, gradeId string, teacher, education *string,
) error {
	stmt := `
		INSERT INTO school.school_year
			(school_year, school_id, grade_id, teacher, education)
		VALUES ($1, $2, $3, $4, $5);`
	_, err := m.DB.Exec(stmt,
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

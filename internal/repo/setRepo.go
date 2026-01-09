package repo

import (
	"database/sql"
	"errors"
	"workout-tracker/internal/models"
)

type SetRepo struct {
	DB *sql.DB
}

func NewSetRepo(db *sql.DB) *SetRepo {
	return &SetRepo{DB: db}
}

func (repo *SetRepo) Create(workoutExerciseId int64, setNumber int, reps *int, weight *float64) (models.Set, error) {
	res, err := repo.DB.Exec(`
		INSERT INTO sets (workout_exercise_id, set_number, reps, weight)
		VALUES (?, ?, ?, ?)
	`, workoutExerciseId, setNumber, reps, weight)
	if err != nil {
		return models.Set{}, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return models.Set{}, err
	}

	return repo.GetById(id)
}

func (repo *SetRepo) GetById(id int64) (models.Set, error) {
	var s models.Set
	err := repo.DB.QueryRow(`
		SELECT id, workout_exercise_id, set_number, reps, weight
		FROM sets
		WHERE id = ?
	`, id).Scan(&s.Id, &s.WorkoutExerciseId, &s.SetNumber, &s.Reps, &s.Weight)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Set{}, ErrNotFound
		}
		return models.Set{}, err
	}
	return s, nil
}

func (repo *SetRepo) ListByWorkoutExercise(workoutExerciseId int64) ([]models.Set, error) {
	rows, err := repo.DB.Query(`
		SELECT id, workout_exercise_id, set_number, reps, weight
		FROM sets
		WHERE workout_exercise_id = ?
		ORDER BY set_number ASC
	`, workoutExerciseId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []models.Set
	for rows.Next() {
		var s models.Set
		if err := rows.Scan(&s.Id, &s.WorkoutExerciseId, &s.SetNumber, &s.Reps, &s.Weight); err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return out, rows.Err()
}

func (repo *SetRepo) Update(id int64, reps *int, weight *float64) (models.Set, error) {
	res, err := repo.DB.Exec(`
		UPDATE sets
		SET reps = ?, weight = ?
		WHERE id = ?
	`, reps, weight, id)
	if err != nil {
		return models.Set{}, err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return models.Set{}, err
	}
	if affected == 0 {
		return models.Set{}, ErrNotFound
	}

	return repo.GetById(id)
}

func (repo *SetRepo) Delete(id int64) error {
	res, err := repo.DB.Exec(`DELETE FROM sets WHERE id = ?`, id)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrNotFound
	}
	return nil
}

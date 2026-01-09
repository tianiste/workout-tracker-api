package repo

import (
	"database/sql"
	"errors"
	"workout-tracker/internal/models"
)

type WorkoutExerciseRepo struct {
	DB *sql.DB
}

func NewWorkoutExerciseRepo(db *sql.DB) *WorkoutExerciseRepo {
	return &WorkoutExerciseRepo{DB: db}
}

func (repo *WorkoutExerciseRepo) AddExercise(workoutId, exerciseId int64, exerciseOrder int, notes *string) (models.WorkoutExercise, error) {
	res, err := repo.DB.Exec(`
		INSERT INTO workout_exercises (workout_id, exercise_id, exercise_order, notes)
		VALUES (?, ?, ?, ?)
	`, workoutId, exerciseId, exerciseOrder, notes)
	if err != nil {
		return models.WorkoutExercise{}, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return models.WorkoutExercise{}, err
	}

	return repo.GetById(id)
}

func (repo *WorkoutExerciseRepo) GetById(id int64) (models.WorkoutExercise, error) {
	var we models.WorkoutExercise
	err := repo.DB.QueryRow(`
		SELECT id, workout_id, exercise_id, exercise_order, notes
		FROM workout_exercises
		WHERE id = ?
	`, id).Scan(&we.Id, &we.WorkoutId, &we.ExerciseId, &we.ExerciseOrder, &we.Notes)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.WorkoutExercise{}, ErrNotFound
		}
		return models.WorkoutExercise{}, err
	}
	return we, nil
}

func (repo *WorkoutExerciseRepo) ListByWorkout(workoutId int64) ([]models.WorkoutExercise, error) {
	rows, err := repo.DB.Query(`
		SELECT id, workout_id, exercise_id, exercise_order, notes
		FROM workout_exercises
		WHERE workout_id = ?
		ORDER BY exercise_order ASC, id ASC
	`, workoutId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []models.WorkoutExercise
	for rows.Next() {
		var we models.WorkoutExercise
		if err := rows.Scan(&we.Id, &we.WorkoutId, &we.ExerciseId, &we.ExerciseOrder, &we.Notes); err != nil {
			return nil, err
		}
		out = append(out, we)
	}
	return out, rows.Err()
}

func (repo *WorkoutExerciseRepo) Update(id int64, exerciseOrder int, notes *string) (models.WorkoutExercise, error) {
	res, err := repo.DB.Exec(`
		UPDATE workout_exercises
		SET exercise_order = ?, notes = ?
		WHERE id = ?
	`, exerciseOrder, notes, id)
	if err != nil {
		return models.WorkoutExercise{}, err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return models.WorkoutExercise{}, err
	}
	if affected == 0 {
		return models.WorkoutExercise{}, ErrNotFound
	}

	return repo.GetById(id)
}

func (repo *WorkoutExerciseRepo) Delete(id int64) error {
	res, err := repo.DB.Exec(`DELETE FROM workout_exercises WHERE id = ?`, id)
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

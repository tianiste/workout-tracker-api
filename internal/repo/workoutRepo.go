package repo

import (
	"database/sql"
	"errors"
	"workout-tracker/internal/models"
)

var ErrNotFound = errors.New("not found")

type WorkoutRepo struct {
	DB *sql.DB
}

func NewWorkoutRepo(db *sql.DB) *WorkoutRepo {
	return &WorkoutRepo{DB: db}
}

func (repo *WorkoutRepo) CreateWorkout(userId int64, performedAt string, durationMinutes *int, notes *string) (models.Workout, error) {
	res, err := repo.DB.Exec(`
		INSERT INTO workouts (user_id, performed_at, duration_minutes, notes)
		VALUES (?, ?, ?, ?)
	`, userId, performedAt, durationMinutes, notes)
	if err != nil {
		return models.Workout{}, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return models.Workout{}, err
	}

	return repo.GetWorkoutById(userId, id)
}

func (repo *WorkoutRepo) GetWorkoutById(userId, workoutId int64) (models.Workout, error) {
	var w models.Workout
	err := repo.DB.QueryRow(`
		SELECT id, user_id, performed_at, duration_minutes, notes, created_at
		FROM workouts
		WHERE id = ? AND user_id = ?
	`, workoutId, userId).Scan(
		&w.Id,
		&w.UserId,
		&w.PerformedAt,
		&w.DurationMinutes,
		&w.Notes,
		&w.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Workout{}, ErrNotFound
		}
		return models.Workout{}, err
	}
	return w, nil
}

func (repo *WorkoutRepo) ListWorkouts(userId int64, limit, offset int) ([]models.Workout, error) {
	rows, err := repo.DB.Query(`
		SELECT id, user_id, performed_at, duration_minutes, notes, created_at
		FROM workouts
		WHERE user_id = ?
		ORDER BY performed_at DESC
		LIMIT ? OFFSET ?
	`, userId, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []models.Workout
	for rows.Next() {
		var w models.Workout
		if err := rows.Scan(&w.Id, &w.UserId, &w.PerformedAt, &w.DurationMinutes, &w.Notes, &w.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, w)
	}
	return out, rows.Err()
}

func (repo *WorkoutRepo) UpdateWorkout(userId, workoutId int64, performedAt string, durationMinutes *int, notes *string) (models.Workout, error) {
	res, err := repo.DB.Exec(`
		UPDATE workouts
		SET performed_at = ?, duration_minutes = ?, notes = ?
		WHERE id = ? AND user_id = ?
	`, performedAt, durationMinutes, notes, workoutId, userId)
	if err != nil {
		return models.Workout{}, err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return models.Workout{}, err
	}
	if affected == 0 {
		return models.Workout{}, ErrNotFound
	}

	return repo.GetWorkoutById(userId, workoutId)
}

func (repo *WorkoutRepo) DeleteWorkout(userId, workoutId int64) error {
	res, err := repo.DB.Exec(`
		DELETE FROM workouts
		WHERE id = ? AND user_id = ?
	`, workoutId, userId)
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

func (repo *WorkoutRepo) MustBeWorkoutOwner(userId, workoutId int64) error {
	var tmp int
	err := repo.DB.QueryRow(`
		SELECT 1
		FROM workouts
		WHERE id = ? AND user_id = ?
	`, workoutId, userId).Scan(&tmp)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNotFound
		}
		return err
	}
	return nil
}

func (repo *WorkoutRepo) GetWorkoutDetails(userId, workoutId int64) (models.WorkoutWithDetails, error) {
	workout, err := repo.GetWorkoutById(userId, workoutId)
	if err != nil {
		return models.WorkoutWithDetails{}, err
	}

	exRows, err := repo.DB.Query(`
		SELECT id, workout_id, exercise_id, exercise_order, notes
		FROM workout_exercises
		WHERE workout_id = ?
		ORDER BY exercise_order ASC, id ASC
	`, workoutId)
	if err != nil {
		return models.WorkoutWithDetails{}, err
	}
	defer exRows.Close()

	var exercises []models.WorkoutExerciseWithSets
	for exRows.Next() {
		var we models.WorkoutExercise
		if err := exRows.Scan(&we.Id, &we.WorkoutId, &we.ExerciseId, &we.ExerciseOrder, &we.Notes); err != nil {
			return models.WorkoutWithDetails{}, err
		}

		sets, err := repo.listSetsByWorkoutExerciseId(we.Id)
		if err != nil {
			return models.WorkoutWithDetails{}, err
		}

		exercises = append(exercises, models.WorkoutExerciseWithSets{
			WorkoutExercise: we,
			Sets:            sets,
		})
	}
	if err := exRows.Err(); err != nil {
		return models.WorkoutWithDetails{}, err
	}

	return models.WorkoutWithDetails{
		Workout:   workout,
		Exercises: exercises,
	}, nil
}

func (repo *WorkoutRepo) listSetsByWorkoutExerciseId(workoutExerciseId int64) ([]models.Set, error) {
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

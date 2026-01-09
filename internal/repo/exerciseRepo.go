package repo

import (
	"database/sql"
	"workout-tracker/internal/models"
)

type ExerciseRepo struct {
	DB *sql.DB
}

func NewExerciseRepo(db *sql.DB) *WorkoutExerciseRepo {
	return &WorkoutExerciseRepo{DB: db}
}

func (repo *ExerciseRepo) ListAllExercises() ([]models.Exercise, error) {
	rows, err := repo.DB.Query(`SELECT id, name, category_id, muscle_group_id FROM exercises`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []models.Exercise
	for rows.Next() {
		var exercise models.Exercise
		if err := rows.Scan(&exercise.Id, &exercise.Name, &exercise.CategoryId, &exercise.MuscleGroupId); err != nil {
			return nil, err
		}
		out = append(out, exercise)
	}
	return out, rows.Err()
}

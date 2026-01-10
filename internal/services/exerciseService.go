package services

import (
	"workout-tracker/internal/models"
	"workout-tracker/internal/repo"
)

type ExerciseService struct {
	Repo *repo.ExerciseRepo
}

func NewExerciseService(repo *repo.ExerciseRepo) *ExerciseService {
	return &ExerciseService{Repo: repo}
}

func (service *ExerciseService) ListAllExercises() ([]models.Exercise, error) {
	out, err := service.Repo.ListAllExercises()
	if err != nil {
		return nil, err
	}
	return out, nil

}

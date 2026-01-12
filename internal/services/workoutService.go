package services

import (
	"errors"
	"fmt"
	"strings"
	"workout-tracker/internal/models"
	"workout-tracker/internal/repo"
)

var ErrForbidden = errors.New("forbidden")

type WorkoutService struct {
	WorkoutRepo         *repo.WorkoutRepo
	WorkoutExerciseRepo *repo.WorkoutExerciseRepo
	SetRepo             *repo.SetRepo
}

func NewWorkoutService(wr *repo.WorkoutRepo, wer *repo.WorkoutExerciseRepo, sr *repo.SetRepo) *WorkoutService {
	return &WorkoutService{
		WorkoutRepo:         wr,
		WorkoutExerciseRepo: wer,
		SetRepo:             sr,
	}
}

func (service *WorkoutService) CreateWorkout(userId int64, performedAt string, durationMinutes *int, notes *string) (models.Workout, error) {
	if strings.TrimSpace(performedAt) == "" {
		return models.Workout{}, fmt.Errorf("performedAt is required")
	}
	return service.WorkoutRepo.CreateWorkout(userId, performedAt, durationMinutes, notes)
}

func (service *WorkoutService) GetWorkout(userId, workoutId int64) (models.Workout, error) {
	return service.WorkoutRepo.GetWorkoutById(userId, workoutId)
}

func (service *WorkoutService) GetWorkoutDetails(userId, workoutId int64) (models.WorkoutWithDetails, error) {
	return service.WorkoutRepo.GetWorkoutDetails(userId, workoutId)
}

func (service *WorkoutService) ListWorkouts(userId int64, limit, offset int) ([]models.Workout, error) {
	if limit <= 0 {
		limit = 25
	}
	if limit > 200 {
		limit = 200
	}
	if offset < 0 {
		offset = 0
	}
	return service.WorkoutRepo.ListWorkouts(userId, limit, offset)
}

func (service *WorkoutService) UpdateWorkout(userId, workoutId int64, performedAt string, durationMinutes *int, notes *string) (models.Workout, error) {
	if strings.TrimSpace(performedAt) == "" {
		return models.Workout{}, fmt.Errorf("performedAt is required")
	}
	return service.WorkoutRepo.UpdateWorkout(userId, workoutId, performedAt, durationMinutes, notes)
}

func (service *WorkoutService) DeleteWorkout(userId, workoutId int64) error {
	return service.WorkoutRepo.DeleteWorkout(userId, workoutId)
}

func (service *WorkoutService) AddExerciseToWorkout(userId, workoutId, exerciseId int64, exerciseOrder int, notes *string) (models.WorkoutExercise, error) {
	if err := service.WorkoutRepo.MustBeWorkoutOwner(userId, workoutId); err != nil {
		return models.WorkoutExercise{}, err
	}
	if exerciseOrder <= 0 {
		return models.WorkoutExercise{}, fmt.Errorf("exerciseOrder must be >= 1")
	}
	return service.WorkoutExerciseRepo.AddExercise(workoutId, exerciseId, exerciseOrder, notes)
}

func (service *WorkoutService) UpdateWorkoutExercise(userId int64, workoutExerciseId int64, exerciseOrder int, notes *string) (models.WorkoutExercise, error) {
	we, err := service.WorkoutExerciseRepo.GetById(workoutExerciseId)
	if err != nil {
		return models.WorkoutExercise{}, err
	}

	if err := service.WorkoutRepo.MustBeWorkoutOwner(userId, we.WorkoutId); err != nil {
		return models.WorkoutExercise{}, err
	}

	if exerciseOrder <= 0 {
		return models.WorkoutExercise{}, fmt.Errorf("exerciseOrder must be >= 1")
	}
	return service.WorkoutExerciseRepo.Update(workoutExerciseId, exerciseOrder, notes)
}

func (service *WorkoutService) DeleteWorkoutExercise(userId int64, workoutExerciseId int64) error {
	we, err := service.WorkoutExerciseRepo.GetById(workoutExerciseId)
	if err != nil {
		return err
	}

	if err := service.WorkoutRepo.MustBeWorkoutOwner(userId, we.WorkoutId); err != nil {
		return err
	}

	return service.WorkoutExerciseRepo.Delete(workoutExerciseId)
}

func (service *WorkoutService) AddSet(userId int64, workoutExerciseId int64, setNumber int, reps *int, weight *float64) (models.Set, error) {
	we, err := service.WorkoutExerciseRepo.GetById(workoutExerciseId)
	if err != nil {
		return models.Set{}, err
	}

	if err := service.WorkoutRepo.MustBeWorkoutOwner(userId, we.WorkoutId); err != nil {
		return models.Set{}, err
	}

	if setNumber <= 0 {
		return models.Set{}, fmt.Errorf("setNumber must be >= 1")
	}

	return service.SetRepo.Create(workoutExerciseId, setNumber, reps, weight)
}

func (service *WorkoutService) UpdateSet(userId int64, setId int64, reps *int, weight *float64) (models.Set, error) {
	set, err := service.SetRepo.GetById(setId)
	if err != nil {
		return models.Set{}, err
	}

	we, err := service.WorkoutExerciseRepo.GetById(set.WorkoutExerciseId)
	if err != nil {
		return models.Set{}, err
	}

	if err := service.WorkoutRepo.MustBeWorkoutOwner(userId, we.WorkoutId); err != nil {
		return models.Set{}, err
	}

	return service.SetRepo.Update(setId, reps, weight)
}

func (service *WorkoutService) DeleteSet(userId int64, setId int64) error {
	set, err := service.SetRepo.GetById(setId)
	if err != nil {
		return err
	}

	we, err := service.WorkoutExerciseRepo.GetById(set.WorkoutExerciseId)
	if err != nil {
		return err
	}

	if err := service.WorkoutRepo.MustBeWorkoutOwner(userId, we.WorkoutId); err != nil {
		return err
	}

	return service.SetRepo.Delete(setId)
}

func (service *WorkoutService) GetWorkoutReport(userId, workoutId int64) (models.WorkoutReport, error) {
	return service.WorkoutRepo.GetWorkoutReport(userId, workoutId)
}

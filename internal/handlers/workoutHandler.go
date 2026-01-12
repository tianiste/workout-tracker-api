package handlers

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"workout-tracker/internal/repo"
	"workout-tracker/internal/services"

	"github.com/gin-gonic/gin"
)

type WorkoutHandler struct {
	Service *services.WorkoutService
}

func NewWorkoutHandler(service *services.WorkoutService) *WorkoutHandler {
	return &WorkoutHandler{Service: service}
}

func getUserIDFromContext(ctx *gin.Context) (int64, error) {
	raw, ok := ctx.Get("userId")
	if !ok {
		return 0, errors.New("missing userId in context")
	}

	if s, ok := raw.(string); ok {
		s = strings.TrimSpace(s)
		if s == "" {
			return 0, errors.New("empty userId in context")
		}
		return strconv.ParseInt(s, 10, 64)
	}

	switch v := raw.(type) {
	case int64:
		return v, nil
	case int:
		return int64(v), nil
	default:
		return 0, errors.New("invalid userId type in context")
	}
}

func parseIDParam(ctx *gin.Context, name string) (int64, error) {
	return strconv.ParseInt(ctx.Param(name), 10, 64)
}

type createWorkoutRequest struct {
	PerformedAt     string  `json:"performedAt" binding:"required"`
	DurationMinutes *int    `json:"durationMinutes"`
	Notes           *string `json:"notes"`
}

type updateWorkoutRequest struct {
	PerformedAt     string  `json:"performedAt" binding:"required"`
	DurationMinutes *int    `json:"durationMinutes"`
	Notes           *string `json:"notes"`
}

type addWorkoutExerciseRequest struct {
	ExerciseId    int64   `json:"exerciseId" binding:"required"`
	ExerciseOrder int     `json:"exerciseOrder" binding:"required"`
	Notes         *string `json:"notes"`
}

type updateWorkoutExerciseRequest struct {
	ExerciseOrder int     `json:"exerciseOrder" binding:"required"`
	Notes         *string `json:"notes"`
}

type addSetRequest struct {
	SetNumber int      `json:"setNumber" binding:"required"`
	Reps      *int     `json:"reps"`
	Weight    *float64 `json:"weight"`
}

type updateSetRequest struct {
	Reps   *int     `json:"reps"`
	Weight *float64 `json:"weight"`
}

func (h *WorkoutHandler) CreateWorkout(ctx *gin.Context) {
	userId, err := getUserIDFromContext(ctx)
	if err != nil {
		log.Printf("[CreateWorkout] unauthorized: %v", err)
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req createWorkoutRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("[CreateWorkout] bad request user=%d: %v", userId, err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	workout, err := h.Service.CreateWorkout(userId, req.PerformedAt, req.DurationMinutes, req.Notes)
	if err != nil {
		log.Printf("[CreateWorkout] failed user=%d: %v", userId, err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, workout)
}

func (h *WorkoutHandler) ListWorkouts(ctx *gin.Context) {
	userId, err := getUserIDFromContext(ctx)
	if err != nil {
		log.Printf("[ListWorkouts] unauthorized: %v", err)
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	limit, limErr := strconv.Atoi(ctx.DefaultQuery("limit", "25"))
	if limErr != nil {
		log.Printf("[ListWorkouts] invalid limit user=%d: %v", userId, limErr)
		limit = 25
	}
	offset, offErr := strconv.Atoi(ctx.DefaultQuery("offset", "0"))
	if offErr != nil {
		log.Printf("[ListWorkouts] invalid offset user=%d: %v", userId, offErr)
		offset = 0
	}

	workouts, err := h.Service.ListWorkouts(userId, limit, offset)
	if err != nil {
		log.Printf("[ListWorkouts] failed user=%d limit=%d offset=%d: %v", userId, limit, offset, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list workouts"})
		return
	}

	ctx.JSON(http.StatusOK, workouts)
}

func (h *WorkoutHandler) GetWorkout(ctx *gin.Context) {
	userId, err := getUserIDFromContext(ctx)
	if err != nil {
		log.Printf("[GetWorkout] unauthorized: %v", err)
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	workoutId, err := parseIDParam(ctx, "id")
	if err != nil {
		log.Printf("[GetWorkout] invalid workout id user=%d: %v", userId, err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid workout id"})
		return
	}

	workout, err := h.Service.GetWorkout(userId, workoutId)
	if err != nil {
		log.Printf("[GetWorkout] failed user=%d workout=%d: %v", userId, workoutId, err)
		if errors.Is(err, repo.ErrNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "workout not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get workout"})
		return
	}

	ctx.JSON(http.StatusOK, workout)
}

func (h *WorkoutHandler) GetWorkoutDetails(ctx *gin.Context) {
	userId, err := getUserIDFromContext(ctx)
	if err != nil {
		log.Printf("[GetWorkoutDetails] unauthorized: %v", err)
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	workoutId, err := parseIDParam(ctx, "id")
	if err != nil {
		log.Printf("[GetWorkoutDetails] invalid workout id user=%d: %v", userId, err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid workout id"})
		return
	}

	details, err := h.Service.GetWorkoutDetails(userId, workoutId)
	if err != nil {
		log.Printf("[GetWorkoutDetails] failed user=%d workout=%d: %v", userId, workoutId, err)
		if errors.Is(err, repo.ErrNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "workout not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get workout details"})
		return
	}

	ctx.JSON(http.StatusOK, details)
}

func (h *WorkoutHandler) UpdateWorkout(ctx *gin.Context) {
	userId, err := getUserIDFromContext(ctx)
	if err != nil {
		log.Printf("[UpdateWorkout] unauthorized: %v", err)
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	workoutId, err := parseIDParam(ctx, "id")
	if err != nil {
		log.Printf("[UpdateWorkout] invalid workout id user=%d: %v", userId, err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid workout id"})
		return
	}

	var req updateWorkoutRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("[UpdateWorkout] bad request user=%d workout=%d: %v", userId, workoutId, err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	updated, err := h.Service.UpdateWorkout(userId, workoutId, req.PerformedAt, req.DurationMinutes, req.Notes)
	if err != nil {
		log.Printf("[UpdateWorkout] failed user=%d workout=%d: %v", userId, workoutId, err)
		if errors.Is(err, repo.ErrNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "workout not found"})
			return
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, updated)
}

func (h *WorkoutHandler) DeleteWorkout(ctx *gin.Context) {
	userId, err := getUserIDFromContext(ctx)
	if err != nil {
		log.Printf("[DeleteWorkout] unauthorized: %v", err)
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	workoutId, err := parseIDParam(ctx, "id")
	if err != nil {
		log.Printf("[DeleteWorkout] invalid workout id user=%d: %v", userId, err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid workout id"})
		return
	}

	if err := h.Service.DeleteWorkout(userId, workoutId); err != nil {
		log.Printf("[DeleteWorkout] failed user=%d workout=%d: %v", userId, workoutId, err)
		if errors.Is(err, repo.ErrNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "workout not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete workout"})
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (h *WorkoutHandler) AddExerciseToWorkout(ctx *gin.Context) {
	userId, err := getUserIDFromContext(ctx)
	if err != nil {
		log.Printf("[AddExerciseToWorkout] unauthorized: %v", err)
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	workoutId, err := parseIDParam(ctx, "id")
	if err != nil {
		log.Printf("[AddExerciseToWorkout] invalid workout id user=%d: %v", userId, err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid workout id"})
		return
	}

	var req addWorkoutExerciseRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("[AddExerciseToWorkout] bad request user=%d workout=%d: %v", userId, workoutId, err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	we, err := h.Service.AddExerciseToWorkout(userId, workoutId, req.ExerciseId, req.ExerciseOrder, req.Notes)
	if err != nil {
		log.Printf("[AddExerciseToWorkout] failed user=%d workout=%d exercise=%d: %v", userId, workoutId, req.ExerciseId, err)
		if errors.Is(err, repo.ErrNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "workout not found"})
			return
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, we)
}

func (h *WorkoutHandler) UpdateWorkoutExercise(ctx *gin.Context) {
	userId, err := getUserIDFromContext(ctx)
	if err != nil {
		log.Printf("[UpdateWorkoutExercise] unauthorized: %v", err)
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	workoutExerciseId, err := parseIDParam(ctx, "id")
	if err != nil {
		log.Printf("[UpdateWorkoutExercise] invalid workout exercise id user=%d: %v", userId, err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid workout exercise id"})
		return
	}

	var req updateWorkoutExerciseRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("[UpdateWorkoutExercise] bad request user=%d workoutExercise=%d: %v", userId, workoutExerciseId, err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	updated, err := h.Service.UpdateWorkoutExercise(userId, workoutExerciseId, req.ExerciseOrder, req.Notes)
	if err != nil {
		log.Printf("[UpdateWorkoutExercise] failed user=%d workoutExercise=%d: %v", userId, workoutExerciseId, err)
		if errors.Is(err, repo.ErrNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "workout exercise not found"})
			return
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, updated)
}

func (h *WorkoutHandler) DeleteWorkoutExercise(ctx *gin.Context) {
	userId, err := getUserIDFromContext(ctx)
	if err != nil {
		log.Printf("[DeleteWorkoutExercise] unauthorized: %v", err)
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	workoutExerciseId, err := parseIDParam(ctx, "id")
	if err != nil {
		log.Printf("[DeleteWorkoutExercise] invalid workout exercise id user=%d: %v", userId, err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid workout exercise id"})
		return
	}

	if err := h.Service.DeleteWorkoutExercise(userId, workoutExerciseId); err != nil {
		log.Printf("[DeleteWorkoutExercise] failed user=%d workoutExercise=%d: %v", userId, workoutExerciseId, err)
		if errors.Is(err, repo.ErrNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "workout exercise not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete workout exercise"})
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (h *WorkoutHandler) AddSet(ctx *gin.Context) {
	userId, err := getUserIDFromContext(ctx)
	if err != nil {
		log.Printf("[AddSet] unauthorized: %v", err)
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	workoutExerciseId, err := parseIDParam(ctx, "id")
	if err != nil {
		log.Printf("[AddSet] invalid workout exercise id user=%d: %v", userId, err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid workout exercise id"})
		return
	}

	var req addSetRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("[AddSet] bad request user=%d workoutExercise=%d: %v", userId, workoutExerciseId, err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	set, err := h.Service.AddSet(userId, workoutExerciseId, req.SetNumber, req.Reps, req.Weight)
	if err != nil {
		log.Printf("[AddSet] failed user=%d workoutExercise=%d setNumber=%d: %v", userId, workoutExerciseId, req.SetNumber, err)
		if errors.Is(err, repo.ErrNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "workout exercise not found"})
			return
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, set)
}

func (h *WorkoutHandler) UpdateSet(ctx *gin.Context) {
	userId, err := getUserIDFromContext(ctx)
	if err != nil {
		log.Printf("[UpdateSet] unauthorized: %v", err)
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	setId, err := parseIDParam(ctx, "id")
	if err != nil {
		log.Printf("[UpdateSet] invalid set id user=%d: %v", userId, err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid set id"})
		return
	}

	var req updateSetRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("[UpdateSet] bad request user=%d set=%d: %v", userId, setId, err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	updated, err := h.Service.UpdateSet(userId, setId, req.Reps, req.Weight)
	if err != nil {
		log.Printf("[UpdateSet] failed user=%d set=%d: %v", userId, setId, err)
		if errors.Is(err, repo.ErrNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "set not found"})
			return
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, updated)
}

func (h *WorkoutHandler) DeleteSet(ctx *gin.Context) {
	userId, err := getUserIDFromContext(ctx)
	if err != nil {
		log.Printf("[DeleteSet] unauthorized: %v", err)
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	setId, err := parseIDParam(ctx, "id")
	if err != nil {
		log.Printf("[DeleteSet] invalid set id user=%d: %v", userId, err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid set id"})
		return
	}

	if err := h.Service.DeleteSet(userId, setId); err != nil {
		log.Printf("[DeleteSet] failed user=%d set=%d: %v", userId, setId, err)
		if errors.Is(err, repo.ErrNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "set not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete set"})
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (h *WorkoutHandler) GetWorkoutReport(ctx *gin.Context) {
	userId, err := getUserIDFromContext(ctx)
	if err != nil {
		log.Printf("[GetWorkoutReport] unauthorized: %v", err)
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	workoutId, err := parseIDParam(ctx, "id")
	if err != nil {
		log.Printf("[GetWorkoutReport] invalid workout id user=%d: %v", userId, err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid workout id"})
		return
	}

	report, err := h.Service.GetWorkoutReport(userId, workoutId)
	if err != nil {
		log.Printf("[GetWorkoutReport] failed user=%d workout=%d: %v", userId, workoutId, err)
		if errors.Is(err, repo.ErrNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "workout not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate workout report"})
		return
	}

	ctx.JSON(http.StatusOK, report)
}

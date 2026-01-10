package handlers

import (
	"log"
	"net/http"
	"workout-tracker/internal/services"

	"github.com/gin-gonic/gin"
)

type ExerciseHandler struct {
	Service *services.ExerciseService
}

func NewExerciseHandler(service *services.ExerciseService) *ExerciseHandler {
	return &ExerciseHandler{Service: service}
}

func (h *ExerciseHandler) ListAllExercises(ctx *gin.Context) {
	exercises, err := h.Service.ListAllExercises()
	if err != nil {
		log.Printf("[ListAllExercises] error %s", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
	ctx.JSON(http.StatusOK, gin.H{"exercises": exercises})
}

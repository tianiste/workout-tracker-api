package handlers

import (
	"log"
	"net/http"
	"workout-tracker/internal/models"

	"github.com/gin-gonic/gin"
)

func Login(ctx *gin.Context) {
	var user models.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		log.Println("incorrect json body", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "missing or invalid json body"})
	}
}

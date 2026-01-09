package handlers

import (
	"log"
	"net/http"

	"workout-tracker/internal/services"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	Service *services.UserService
}

func NewUserHandler(service *services.UserService) *UserHandler {
	return &UserHandler{Service: service}
}

type authRequest struct {
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *UserHandler) Login(ctx *gin.Context) {
	var req authRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("[LOGIN] invalid request body: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	token, err := h.Service.Login(req.Name, req.Password)
	if err != nil {
		log.Printf("[LOGIN] failed for user=%q: %v", req.Name, err)
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"token": token})
}

func (h *UserHandler) Register(ctx *gin.Context) {
	var req authRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("[REGISTER] invalid request body: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	userId, err := h.Service.Register(req.Name, req.Password)
	if err != nil {
		log.Printf("[REGISTER] failed for user=%q: %v", req.Name, err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"userId": userId})
}

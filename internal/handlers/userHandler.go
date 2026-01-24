package handlers

import (
	"log"
	"net/http"
	"time"

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

func setRefreshCookie(ctx *gin.Context, rawToken string, expiresAt time.Time) {
	c := &http.Cookie{
		Name:     "refresh_token",
		Value:    rawToken,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Path:     "/refresh",
		Expires:  expiresAt,
	}
	http.SetCookie(ctx.Writer, c)
}

func clearRefreshCookie(ctx *gin.Context) {
	c := &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Path:     "/refresh",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
	}
	http.SetCookie(ctx.Writer, c)
}

func (h *UserHandler) Login(ctx *gin.Context) {
	var req authRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("[LOGIN] invalid request body: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	token, rawToken, expiresAt, err := h.Service.Login(ctx.Request.Context(), req.Name, req.Password)
	if err != nil {
		log.Printf("[LOGIN] failed for user=%q: %v", req.Name, err)
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}
	setRefreshCookie(ctx, rawToken, expiresAt)

	ctx.JSON(http.StatusOK, gin.H{"token": token})
}

func (h *UserHandler) Refresh(ctx *gin.Context) {
	rawToken, err := ctx.Cookie("refresh_token")
	if err != nil || rawToken == "" {
		clearRefreshCookie(ctx)
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "missing refresh token"})
		return
	}

	newAccess, newRefresh, newRefreshExp, err := h.Service.Refresh(ctx.Request.Context(), rawToken)
	if err != nil {
		if err == services.ErrMissingRefreshToken || err == services.ErrInvalidRefreshToken {
			clearRefreshCookie(ctx)
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
			return
		}
		log.Printf("[REFRESH] server error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
		return
	}

	setRefreshCookie(ctx, newRefresh, newRefreshExp)

	ctx.JSON(http.StatusOK, gin.H{"token": newAccess})
}

func (h *UserHandler) Logout(ctx *gin.Context) {
	rawToken, _ := ctx.Cookie("refresh_token")

	clearRefreshCookie(ctx)

	if rawToken == "" {
		ctx.Status(http.StatusNoContent)
		return
	}

	err := h.Service.Logout(ctx.Request.Context(), rawToken)
	if err != nil && err != services.ErrMissingRefreshToken && err != services.ErrInvalidRefreshToken {
		log.Printf("[LOGOUT] server error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
		return
	}

	ctx.Status(http.StatusNoContent)
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

package middleware

import (
	"net/http"
	"strings"

	"workout-tracker/internal/services"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(userService *services.UserService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		auth := ctx.GetHeader("Authorization")
		if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
			ctx.Abort()
			return
		}

		tokenString := strings.TrimSpace(strings.TrimPrefix(auth, "Bearer "))
		if tokenString == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
			ctx.Abort()
			return
		}

		claims, err := userService.VerifyToken(tokenString)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			ctx.Abort()
			return
		}

		sub, ok := claims["sub"].(string)
		if !ok || sub == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			ctx.Abort()
			return
		}

		ctx.Set("userId", sub)
		ctx.Next()
	}
}

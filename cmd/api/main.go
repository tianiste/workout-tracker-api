package main

import (
	"log"
	"net/http"
	"time"

	ratelimit "github.com/JGLTechnologies/gin-rate-limit"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"workout-tracker/db"
	"workout-tracker/internal/handlers"
	"workout-tracker/internal/middleware"
	"workout-tracker/internal/repo"
	"workout-tracker/internal/services"
)

func keyFunc(c *gin.Context) string {
	return c.ClientIP()
}

func errorHandler(c *gin.Context, info ratelimit.Info) {
	c.String(429, "Too many requests. Try again in "+time.Until(info.ResetTime).String())
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("warning: failed to load .env: %v", err)
	}

	db.InitDB()
	defer func() {
		if db.DB != nil {
			_ = db.DB.Close()
		}
	}()

	router := gin.Default()

	store := ratelimit.InMemoryStore(&ratelimit.InMemoryOptions{
		Rate:  time.Second,
		Limit: 5,
	})
	rateLimiter := ratelimit.RateLimiter(store, &ratelimit.Options{
		ErrorHandler: errorHandler,
		KeyFunc:      keyFunc,
	})

	userRepo := repo.NewUserRepo(db.DB)
	userService := services.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)

	exerciseRepo := repo.NewExerciseRepo(db.DB)
	exerciseService := services.NewExerciseService(exerciseRepo)
	exerciseHandler := handlers.NewExerciseHandler(exerciseService)

	workoutRepo := repo.NewWorkoutRepo(db.DB)
	workoutExerciseRepo := repo.NewWorkoutExerciseRepo(db.DB)
	setRepo := repo.NewSetRepo(db.DB)
	workoutService := services.NewWorkoutService(workoutRepo, workoutExerciseRepo, setRepo)
	workoutHandler := handlers.NewWorkoutHandler(workoutService)

	api := router.Group("/api")
	api.Use(rateLimiter)
	{
		api.GET("/ping", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{"message": "pong"})
		})
		api.POST("/register", userHandler.Register)
		api.POST("/login", userHandler.Login)

		authorized := api.Group("/")
		authorized.Use(middleware.AuthMiddleware(userService))
		{
			authorized.POST("/workouts", workoutHandler.CreateWorkout)
			authorized.GET("/workouts", workoutHandler.ListWorkouts)
			authorized.GET("/workouts/:id", workoutHandler.GetWorkout)
			authorized.GET("/workouts/:id/details", workoutHandler.GetWorkoutDetails)
			authorized.PUT("/workouts/:id", workoutHandler.UpdateWorkout)
			authorized.DELETE("/workouts/:id", workoutHandler.DeleteWorkout)

			authorized.GET("/exercises", exerciseHandler.ListAllExercises)

			authorized.POST("/workouts/:id/exercises", workoutHandler.AddExerciseToWorkout)
			authorized.PUT("/workout-exercises/:id", workoutHandler.UpdateWorkoutExercise)
			authorized.DELETE("/workout-exercises/:id", workoutHandler.DeleteWorkoutExercise)

			authorized.POST("/workout-exercises/:id/sets", workoutHandler.AddSet)
			authorized.PUT("/sets/:id", workoutHandler.UpdateSet)
			authorized.DELETE("/sets/:id", workoutHandler.DeleteSet)
		}
	}

	if err := router.Run(); err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}

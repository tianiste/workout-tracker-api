package models

type Workout struct {
	Id              int64   `json:"id"`
	UserId          int64   `json:"userId"`
	PerformedAt     string  `json:"performedAt"`
	DurationMinutes *int    `json:"durationMinutes,omitempty"`
	Notes           *string `json:"notes,omitempty"`
	CreatedAt       string  `json:"createdAt"`
}

type WorkoutExercise struct {
	Id            int64   `json:"id"`
	WorkoutId     int64   `json:"workoutId"`
	ExerciseId    int64   `json:"exerciseId"`
	ExerciseOrder int     `json:"exerciseOrder"`
	Notes         *string `json:"notes,omitempty"`
}

type Set struct {
	Id                int64    `json:"id"`
	WorkoutExerciseId int64    `json:"workoutExerciseId"`
	SetNumber         int      `json:"setNumber"`
	Reps              *int     `json:"reps,omitempty"`
	Weight            *float64 `json:"weight,omitempty"`
}

type WorkoutExerciseWithSets struct {
	WorkoutExercise
	Sets []Set `json:"sets"`
}

type WorkoutWithDetails struct {
	Workout
	Exercises []WorkoutExerciseWithSets `json:"exercises"`
}

package models

type WorkoutReportExercise struct {
	ExerciseId   int64    `json:"exerciseId"`
	ExerciseName string   `json:"exerciseName"`
	SetsCount    int      `json:"setsCount"`
	TotalReps    int      `json:"totalReps"`
	MaxWeight    *float64 `json:"maxWeight,omitempty"`
	TotalVolume  float64  `json:"totalVolume"`
}

type WorkoutReport struct {
	WorkoutId       int64   `json:"workoutId"`
	UserId          int64   `json:"userId"`
	PerformedAt     string  `json:"performedAt"`
	DurationMinutes *int    `json:"durationMinutes,omitempty"`
	Notes           *string `json:"notes,omitempty"`
	CreatedAt       string  `json:"createdAt"`

	TotalExercises int     `json:"totalExercises"`
	TotalSets      int     `json:"totalSets"`
	TotalReps      int     `json:"totalReps"`
	TotalVolume    float64 `json:"totalVolume"`

	Exercises []WorkoutReportExercise `json:"exercises"`
}

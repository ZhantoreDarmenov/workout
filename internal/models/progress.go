package models

import "time"

// DayDetails includes a day's exercises and food
// to present workout instructions to a client.
type DayDetails struct {
	Day      Days      `json:"day"`
	Food     Food      `json:"food"`
	Exercise Exercises `json:"exercise"`
}

// ProgramProgress represents completion of a workout day by a client.
type ProgramProgress struct {
	ID                int        `json:"id"`
	ClientID          int        `json:"client_id"`
	DayID             int        `json:"day_id"`
	FoodCompleted     bool       `json:"food_completed"`
	ExerciseCompleted bool       `json:"exercise_completed"`
	Completed         *time.Time `json:"completed,omitempty"`
}

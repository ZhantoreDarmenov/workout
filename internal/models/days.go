package models

import (
	"time"
)

type Days struct {
	ID               int        `json:"id"`
	WorkOutProgramID int        `json:"work_out_program_id"`
	ExercisesID      int        `json:"exercises_id"`
	FoodID           int        `json:"food_id"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        *time.Time `json:"updated_at,omitempty"`
}

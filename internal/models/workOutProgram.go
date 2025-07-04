package models

import (
	"time"
)

type WorkOutProgram struct {
	ID          int        `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Sets        string     `json:"sets"`
	Repetitions int        `json:"repetitions"`
	Exercises   string     `json:"exercises"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
}

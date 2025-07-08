package models

import (
	"time"
)

type WorkOutProgram struct {
	ID          int        `json:"id"`
	TrainerID   int        `json:"trainer_id"`
	Name        string     `json:"name"`
	Days        int        `json:"days"`
	Description string     `json:"description"`
	Duration    string     `json:"duration,omitempty"`
	Clients     string     `json:"clients,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
}

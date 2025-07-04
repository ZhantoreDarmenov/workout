package models

import (
	"time"
)

type TrainerClient struct {
	ID        int        `json:"id"`
	ClientID  string     `json:"client_id"`
	TrainerID string     `json:"trainer_id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

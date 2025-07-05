package models

import (
	"time"
)

type Exercises struct {
	ID          int        `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Sets        string     `json:"sets"`
	MediaURL    string     `json:"media_url"`
	Repetitions string     `json:"repetitions"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
}

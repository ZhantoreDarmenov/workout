package models

import "time"

// ProgramInvite represents an invitation for a client to join a workout program.
type ProgramInvite struct {
	ID            int        `json:"id"`
	ProgramID     int        `json:"program_id"`
	Email         string     `json:"email"`
	Message       string     `json:"message"`
	AccessDays    int        `json:"access_days"`
	Token         string     `json:"token"`
	ClientID      *int       `json:"client_id,omitempty"`
	AcceptedAt    *time.Time `json:"accepted_at,omitempty"`
	AccessExpires *time.Time `json:"access_expires,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     *time.Time `json:"updated_at,omitempty"`
}

package models

// ClientProgress describes progress statistics for a client.
type ClientProgress struct {
	ClientID      int     `json:"client_id"`
	CompletedDays int     `json:"completed_days"`
	TotalDays     int     `json:"total_days"`
	Progress      float64 `json:"progress"`
}

// TrainerAnalytics aggregates program and client statistics for a trainer.
type TrainerAnalytics struct {
	ProgramCount     int              `json:"program_count"`
	ClientCount      int              `json:"client_count"`
	CompletedClients int              `json:"completed_clients"`
	AverageProgress  float64          `json:"average_progress"`
	ClientsProgress  []ClientProgress `json:"clients_progress"`
}

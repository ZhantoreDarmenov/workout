package models

import (
	"time"
)

type Food struct {
	ID            int        `json:"id"`
	Name          string     `json:"name"`
	Description   string     `json:"description"`
	Calories      float64    `json:"calories"`
	Protein       float64    `json:"protein"`
	Fats          float64    `json:"fats"`
	Carbohydrates float64    `json:"carbohydrates"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     *time.Time `json:"updated_at,omitempty"`
}

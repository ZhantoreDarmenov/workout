package repositories

import (
	"context"
	"database/sql"
	"time"

	"workout/internal/models"
)

// ExerciseRepository handles CRUD for exercises.
type ExerciseRepository struct {
	DB *sql.DB
}

func (r *ExerciseRepository) CreateExercise(ctx context.Context, ex models.Exercises) (models.Exercises, error) {
	query := `INSERT INTO exercises (name, description, media_url, sets, repetitions, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)`
	ex.CreatedAt = time.Now()
	ex.UpdatedAt = &ex.CreatedAt
	res, err := r.DB.ExecContext(ctx, query, ex.Name, ex.Description, ex.MediaURL, ex.Sets, ex.Repetitions, ex.CreatedAt, ex.UpdatedAt)
	if err != nil {
		return models.Exercises{}, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return models.Exercises{}, err
	}
	ex.ID = int(id)
	return ex, nil
}

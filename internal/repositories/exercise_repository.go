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

// UpdateExercise updates an exercise by ID.
func (r *ExerciseRepository) UpdateExercise(ctx context.Context, ex models.Exercises) (models.Exercises, error) {
	now := time.Now()
	ex.UpdatedAt = &now
	query := `UPDATE exercises SET name = ?, description = ?, media_url = ?, sets = ?, repetitions = ?, updated_at = ? WHERE id = ?`
	res, err := r.DB.ExecContext(ctx, query, ex.Name, ex.Description, ex.MediaURL, ex.Sets, ex.Repetitions, ex.UpdatedAt, ex.ID)
	if err != nil {
		return models.Exercises{}, err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return models.Exercises{}, err
	}
	if rows == 0 {
		return models.Exercises{}, models.ErrExerciseNotFound
	}
	return ex, nil
}

// DeleteExercise removes an exercise by ID.
func (r *ExerciseRepository) DeleteExercise(ctx context.Context, id int) error {
	res, err := r.DB.ExecContext(ctx, `DELETE FROM exercises WHERE id = ?`, id)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return models.ErrExerciseNotFound
	}
	return nil
}

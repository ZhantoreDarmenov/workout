package repositories

import (
	"context"
	"database/sql"
	"time"

	"workout/internal/models"
)

type ProgramRepository struct {
	DB *sql.DB
}

func (r *ProgramRepository) CreateProgram(ctx context.Context, p models.WorkOutProgram) (models.WorkOutProgram, error) {
	query := `
INSERT INTO workout_programs (trainer_id, name, days, description, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?)
`
	p.CreatedAt = time.Now()
	p.UpdatedAt = &p.CreatedAt
	res, err := r.DB.ExecContext(ctx, query, p.TrainerID, p.Name, p.Days, p.Description, p.CreatedAt, p.UpdatedAt)
	if err != nil {
		return models.WorkOutProgram{}, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return models.WorkOutProgram{}, err
	}
	p.ID = int(id)
	return p, nil
}

func (r *ProgramRepository) GetProgramsByTrainer(ctx context.Context, trainerID int) ([]models.WorkOutProgram, error) {
	query := `
SELECT id, trainer_id, name, days, description, created_at, updated_at
FROM workout_programs
WHERE trainer_id = ?
`
	rows, err := r.DB.QueryContext(ctx, query, trainerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	programs := []models.WorkOutProgram{}
	for rows.Next() {
		var p models.WorkOutProgram
		if err := rows.Scan(&p.ID, &p.TrainerID, &p.Name, &p.Days, &p.Description, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		programs = append(programs, p)
	}
	return programs, rows.Err()
}

// GetProgramByID fetches a workout program by its ID.
func (r *ProgramRepository) GetProgramByID(ctx context.Context, id int) (models.WorkOutProgram, error) {
	var p models.WorkOutProgram
	query := `SELECT id, trainer_id, name, days, description, created_at, updated_at FROM workout_programs WHERE id = ?`
	err := r.DB.QueryRowContext(ctx, query, id).Scan(&p.ID, &p.TrainerID, &p.Name, &p.Days, &p.Description, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.WorkOutProgram{}, models.ErrWorkoutProgramNotFound
		}
		return models.WorkOutProgram{}, err
	}
	return p, nil
}

// UpdateProgram updates an existing workout program.
func (r *ProgramRepository) UpdateProgram(ctx context.Context, p models.WorkOutProgram) (models.WorkOutProgram, error) {
	now := time.Now()
	p.UpdatedAt = &now
	query := `UPDATE workout_programs SET name = ?, days = ?, description = ?, updated_at = ? WHERE id = ?`
	res, err := r.DB.ExecContext(ctx, query, p.Name, p.Days, p.Description, p.UpdatedAt, p.ID)
	if err != nil {
		return models.WorkOutProgram{}, err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return models.WorkOutProgram{}, err
	}
	if rows == 0 {
		return models.WorkOutProgram{}, models.ErrWorkoutProgramNotFound
	}
	return p, nil
}

// DeleteProgram removes a workout program and all of its days.
func (r *ProgramRepository) DeleteProgram(ctx context.Context, id int) error {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	if _, err = tx.ExecContext(ctx, `DELETE FROM days WHERE work_out_program_id = ?`, id); err != nil {
		tx.Rollback()
		return err
	}

	res, err := tx.ExecContext(ctx, `DELETE FROM workout_programs WHERE id = ?`, id)
	if err != nil {
		tx.Rollback()
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		tx.Rollback()
		return err
	}
	if rows == 0 {
		tx.Rollback()
		return models.ErrWorkoutProgramNotFound
	}

	return tx.Commit()
}

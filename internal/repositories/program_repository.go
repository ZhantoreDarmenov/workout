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

	var programs []models.WorkOutProgram
	for rows.Next() {
		var p models.WorkOutProgram
		if err := rows.Scan(&p.ID, &p.TrainerID, &p.Name, &p.Days, &p.Description, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		programs = append(programs, p)
	}
	return programs, rows.Err()
}

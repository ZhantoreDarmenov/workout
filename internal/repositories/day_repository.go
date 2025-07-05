package repositories

import (
	"context"
	"database/sql"
	"time"

	"workout/internal/models"
)

// DayRepository handles workout day and progress queries.
type DayRepository struct {
	DB *sql.DB
}

func (r *DayRepository) GetDayDetails(ctx context.Context, programID, dayNumber int) (models.DayDetails, error) {
	var d models.Days
	err := r.DB.QueryRowContext(ctx, `SELECT id, work_out_program_id, day_number, exercises_id, food_id, created_at, updated_at FROM days WHERE work_out_program_id=? AND day_number=?`, programID, dayNumber).Scan(&d.ID, &d.WorkOutProgramID, &d.DayNumber, &d.ExercisesID, &d.FoodID, &d.CreatedAt, &d.UpdatedAt)
	if err != nil {
		return models.DayDetails{}, err
	}

	var ex models.Exercises
	err = r.DB.QueryRowContext(ctx, `SELECT id, name, description, sets, repetitions, created_at, updated_at FROM exercises WHERE id=?`, d.ExercisesID).Scan(&ex.ID, &ex.Name, &ex.Description, &ex.Sets, &ex.Repetitions, &ex.CreatedAt, &ex.UpdatedAt)
	if err != nil {
		return models.DayDetails{}, err
	}

	var food models.Food
	err = r.DB.QueryRowContext(ctx, `SELECT id, name, description, calories, protein, fats, carbohydrates, created_at, updated_at FROM food WHERE id=?`, d.FoodID).Scan(&food.ID, &food.Name, &food.Description, &food.Calories, &food.Protein, &food.Fats, &food.Carbohydrates, &food.CreatedAt, &food.UpdatedAt)
	if err != nil {
		return models.DayDetails{}, err
	}

	return models.DayDetails{Day: d, Food: food, Exercise: ex}, nil
}

func (r *DayRepository) MarkDayCompleted(ctx context.Context, clientID, dayID int) (models.ProgramProgress, error) {
	prog := models.ProgramProgress{ClientID: clientID, DayID: dayID, Completed: time.Now()}
	res, err := r.DB.ExecContext(ctx, `INSERT INTO progress (client_id, day_id, completed) VALUES (?, ?, ?)`, prog.ClientID, prog.DayID, prog.Completed)
	if err != nil {
		return models.ProgramProgress{}, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return models.ProgramProgress{}, err
	}
	prog.ID = int(id)
	return prog, nil
}

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

func (r *DayRepository) CreateDay(ctx context.Context, day models.Days) (models.Days, error) {
	query := `INSERT INTO days (work_out_program_id, day_number, exercises_id, food_id, created_at, updated_at)
                  VALUES (?, ?, ?, ?, ?, ?)`
	day.CreatedAt = time.Now()
	day.UpdatedAt = &day.CreatedAt
	res, err := r.DB.ExecContext(ctx, query, day.WorkOutProgramID, day.DayNumber, day.ExercisesID, day.FoodID, day.CreatedAt, day.UpdatedAt)
	if err != nil {
		return models.Days{}, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return models.Days{}, err
	}
	day.ID = int(id)
	return day, nil
}

func (r *DayRepository) DaysByProgram(ctx context.Context, programID int) ([]models.DayDetails, error) {
	rows, err := r.DB.QueryContext(ctx, `SELECT d.id, d.work_out_program_id, d.day_number, d.exercises_id, d.food_id,
                d.created_at, d.updated_at,
                e.id, e.name, e.description, e.media_url, e.sets, e.repetitions, e.created_at, e.updated_at,
                f.id, f.name, f.description, f.calories, f.protein, f.fats, f.carbohydrates, f.created_at, f.updated_at
                FROM days d
                JOIN exercises e ON d.exercises_id = e.id
                JOIN food f ON d.food_id = f.id
                WHERE d.work_out_program_id = ? ORDER BY d.day_number`, programID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.DayDetails
	for rows.Next() {
		var d models.Days
		var ex models.Exercises
		var food models.Food
		err = rows.Scan(&d.ID, &d.WorkOutProgramID, &d.DayNumber, &d.ExercisesID, &d.FoodID,
			&d.CreatedAt, &d.UpdatedAt,
			&ex.ID, &ex.Name, &ex.Description, &ex.MediaURL, &ex.Sets, &ex.Repetitions, &ex.CreatedAt, &ex.UpdatedAt,
			&food.ID, &food.Name, &food.Description, &food.Calories, &food.Protein, &food.Fats, &food.Carbohydrates, &food.CreatedAt, &food.UpdatedAt)
		if err != nil {
			return nil, err
		}
		result = append(result, models.DayDetails{Day: d, Exercise: ex, Food: food})
	}
	return result, rows.Err()
}

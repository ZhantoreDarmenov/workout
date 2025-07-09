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
	err := r.DB.QueryRowContext(ctx, `SELECT id, work_out_program_id, day_number, exercises_id, food_id, note, created_at, updated_at FROM days WHERE work_out_program_id=? AND day_number=?`, programID, dayNumber).Scan(&d.ID, &d.WorkOutProgramID, &d.DayNumber, &d.ExercisesID, &d.FoodID, &d.Note, &d.CreatedAt, &d.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.DayDetails{}, models.ErrDayNotFound
		}
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
	now := time.Now()
	// try update existing progress record
	res, err := r.DB.ExecContext(ctx, `UPDATE progress SET completed = ? WHERE client_id = ? AND day_id = ?`, now, clientID, dayID)
	if err != nil {
		return models.ProgramProgress{}, err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return models.ProgramProgress{}, err
	}
	if rows == 0 {
		// insert new progress if none exists
		res, err = r.DB.ExecContext(ctx, `INSERT INTO progress (client_id, day_id, food_completed, exercise_completed, completed) VALUES (?, ?, false, false, ?)`, clientID, dayID, now)
		if err != nil {
			return models.ProgramProgress{}, err
		}
	}

	var prog models.ProgramProgress
	var completed sql.NullTime
	err = r.DB.QueryRowContext(ctx, `SELECT id, client_id, day_id, food_completed, exercise_completed, completed FROM progress WHERE client_id = ? AND day_id = ?`, clientID, dayID).Scan(
		&prog.ID, &prog.ClientID, &prog.DayID, &prog.FoodCompleted, &prog.ExerciseCompleted, &completed)
	if err != nil {
		return models.ProgramProgress{}, err
	}
	if completed.Valid {
		prog.Completed = &completed.Time
	}
	return prog, nil
}

func (r *DayRepository) MarkFoodCompleted(ctx context.Context, clientID, dayID int) (models.ProgramProgress, error) {
	res, err := r.DB.ExecContext(ctx, `UPDATE progress SET food_completed = TRUE WHERE client_id = ? AND day_id = ?`, clientID, dayID)
	if err != nil {
		return models.ProgramProgress{}, err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return models.ProgramProgress{}, err
	}
	if rows == 0 {
		res, err = r.DB.ExecContext(ctx, `INSERT INTO progress (client_id, day_id, food_completed, exercise_completed) VALUES (?, ?, TRUE, FALSE)`, clientID, dayID)
		if err != nil {
			return models.ProgramProgress{}, err
		}
	}
	return r.GetProgress(ctx, clientID, dayID)
}

func (r *DayRepository) MarkExerciseCompleted(ctx context.Context, clientID, dayID int) (models.ProgramProgress, error) {
	res, err := r.DB.ExecContext(ctx, `UPDATE progress SET exercise_completed = TRUE WHERE client_id = ? AND day_id = ?`, clientID, dayID)
	if err != nil {
		return models.ProgramProgress{}, err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return models.ProgramProgress{}, err
	}
	if rows == 0 {
		res, err = r.DB.ExecContext(ctx, `INSERT INTO progress (client_id, day_id, food_completed, exercise_completed) VALUES (?, ?, FALSE, TRUE)`, clientID, dayID)
		if err != nil {
			return models.ProgramProgress{}, err
		}
	}
	return r.GetProgress(ctx, clientID, dayID)
}

func (r *DayRepository) GetProgress(ctx context.Context, clientID, dayID int) (models.ProgramProgress, error) {
	var prog models.ProgramProgress
	var completed sql.NullTime
	err := r.DB.QueryRowContext(ctx, `SELECT id, client_id, day_id, food_completed, exercise_completed, completed FROM progress WHERE client_id = ? AND day_id = ?`, clientID, dayID).Scan(
		&prog.ID, &prog.ClientID, &prog.DayID, &prog.FoodCompleted, &prog.ExerciseCompleted, &completed)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.ProgramProgress{}, models.ErrDayNotFound
		}
		return models.ProgramProgress{}, err
	}
	if completed.Valid {
		prog.Completed = &completed.Time
	}
	return prog, nil
}

// GetProgramProgress returns progress info for all days in a program for a client.
func (r *DayRepository) GetProgramProgress(ctx context.Context, clientID, programID int) ([]models.DayProgressStatus, error) {
	rows, err := r.DB.QueryContext(ctx, `SELECT d.id, d.day_number,
        COALESCE(p.food_completed, FALSE),
        COALESCE(p.exercise_completed, FALSE),
        p.completed
        FROM days d
        LEFT JOIN progress p ON p.day_id = d.id AND p.client_id = ?
        WHERE d.work_out_program_id = ?
        ORDER BY d.day_number`, clientID, programID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.DayProgressStatus
	for rows.Next() {
		var dp models.DayProgressStatus
		var completed sql.NullTime
		if err := rows.Scan(&dp.DayID, &dp.DayNumber, &dp.FoodCompleted, &dp.ExerciseCompleted, &completed); err != nil {
			return nil, err
		}
		if completed.Valid {
			dp.Completed = &completed.Time
		}
		result = append(result, dp)
	}
	return result, rows.Err()
}

func (r *DayRepository) CreateDay(ctx context.Context, day models.Days) (models.Days, error) {
	// ensure referenced records exist to avoid foreign key errors
	var exists bool
	if err := r.DB.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM workout_programs WHERE id = ?)", day.WorkOutProgramID).Scan(&exists); err != nil {
		return models.Days{}, err
	}
	if !exists {
		return models.Days{}, models.ErrWorkoutProgramNotFound
	}

	if err := r.DB.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM exercises WHERE id = ?)", day.ExercisesID).Scan(&exists); err != nil {
		return models.Days{}, err
	}
	if !exists {
		return models.Days{}, models.ErrExerciseNotFound
	}

	if err := r.DB.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM food WHERE id = ?)", day.FoodID).Scan(&exists); err != nil {
		return models.Days{}, err
	}
	if !exists {
		return models.Days{}, models.ErrFoodNotFound
	}

	query := `INSERT INTO days (work_out_program_id, day_number, exercises_id, food_id, note, created_at, updated_at)
                  VALUES (?, ?, ?, ?, ?, ?, ?)`
	day.CreatedAt = time.Now()
	day.UpdatedAt = &day.CreatedAt
	res, err := r.DB.ExecContext(ctx, query, day.WorkOutProgramID, day.DayNumber, day.ExercisesID, day.FoodID, day.Note, day.CreatedAt, day.UpdatedAt)
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
	rows, err := r.DB.QueryContext(ctx, `SELECT d.id, d.work_out_program_id, d.day_number, d.exercises_id, d.food_id, d.note,
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

	result := []models.DayDetails{}
	for rows.Next() {
		var d models.Days
		var ex models.Exercises
		var food models.Food
		err = rows.Scan(&d.ID, &d.WorkOutProgramID, &d.DayNumber, &d.ExercisesID, &d.FoodID, &d.Note,
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

// UpdateDay updates a workout day by its ID.
func (r *DayRepository) UpdateDay(ctx context.Context, day models.Days) (models.Days, error) {
	// ensure referenced records exist to avoid foreign key violations
	var exists bool
	if err := r.DB.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM workout_programs WHERE id = ?)", day.WorkOutProgramID).Scan(&exists); err != nil {
		return models.Days{}, err
	}
	if !exists {
		return models.Days{}, models.ErrWorkoutProgramNotFound
	}

	if err := r.DB.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM exercises WHERE id = ?)", day.ExercisesID).Scan(&exists); err != nil {
		return models.Days{}, err
	}
	if !exists {
		return models.Days{}, models.ErrExerciseNotFound
	}

	if err := r.DB.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM food WHERE id = ?)", day.FoodID).Scan(&exists); err != nil {
		return models.Days{}, err
	}
	if !exists {
		return models.Days{}, models.ErrFoodNotFound
	}

	now := time.Now()
	day.UpdatedAt = &now
	res, err := r.DB.ExecContext(ctx, `UPDATE days SET work_out_program_id = ?, day_number = ?, exercises_id = ?, food_id = ?, note = ?, updated_at = ? WHERE id = ?`,
		day.WorkOutProgramID, day.DayNumber, day.ExercisesID, day.FoodID, day.Note, day.UpdatedAt, day.ID)
	if err != nil {
		return models.Days{}, err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return models.Days{}, err
	}
	if rows == 0 {
		return models.Days{}, models.ErrDayNotFound
	}
	return day, nil
}

// DeleteDay removes a workout day by its ID.
func (r *DayRepository) DeleteDay(ctx context.Context, id int) error {
	res, err := r.DB.ExecContext(ctx, `DELETE FROM days WHERE id = ?`, id)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return models.ErrDayNotFound
	}
	return nil
}

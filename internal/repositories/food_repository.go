package repositories

import (
	"context"
	"database/sql"
	"time"

	"workout/internal/models"
)

// FoodRepository handles CRUD operations for food.
type FoodRepository struct {
	DB *sql.DB
}

func (r *FoodRepository) CreateFood(ctx context.Context, f models.Food) (models.Food, error) {
	query := `INSERT INTO food (name, description, calories, protein, fats, carbohydrates, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	f.CreatedAt = time.Now()
	f.UpdatedAt = &f.CreatedAt
	res, err := r.DB.ExecContext(ctx, query, f.Name, f.Description, f.Calories, f.Protein, f.Fats, f.Carbohydrates, f.CreatedAt, f.UpdatedAt)
	if err != nil {
		return models.Food{}, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return models.Food{}, err
	}
	f.ID = int(id)
	return f, nil
}

// UpdateFood updates food entry by ID.
func (r *FoodRepository) UpdateFood(ctx context.Context, f models.Food) (models.Food, error) {
	now := time.Now()
	f.UpdatedAt = &now
	query := `UPDATE food SET name = ?, description = ?, calories = ?, protein = ?, fats = ?, carbohydrates = ?, updated_at = ? WHERE id = ?`
	res, err := r.DB.ExecContext(ctx, query, f.Name, f.Description, f.Calories, f.Protein, f.Fats, f.Carbohydrates, f.UpdatedAt, f.ID)
	if err != nil {
		return models.Food{}, err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return models.Food{}, err
	}
	if rows == 0 {
		return models.Food{}, models.ErrFoodNotFound
	}
	return f, nil
}

// DeleteFood removes food entry by ID.
func (r *FoodRepository) DeleteFood(ctx context.Context, id int) error {
	res, err := r.DB.ExecContext(ctx, `DELETE FROM food WHERE id = ?`, id)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return models.ErrFoodNotFound
	}
	return nil
}

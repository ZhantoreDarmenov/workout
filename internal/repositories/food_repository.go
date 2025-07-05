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

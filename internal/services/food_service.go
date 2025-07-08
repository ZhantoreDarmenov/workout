package services

import (
	"context"

	"workout/internal/models"
	"workout/internal/repositories"
)

// FoodService contains business logic for food items.
type FoodService struct {
	Repo *repositories.FoodRepository
}

func (s *FoodService) CreateFood(ctx context.Context, f models.Food) (models.Food, error) {
	return s.Repo.CreateFood(ctx, f)
}

func (s *FoodService) UpdateFood(ctx context.Context, f models.Food) (models.Food, error) {
	return s.Repo.UpdateFood(ctx, f)
}

func (s *FoodService) DeleteFood(ctx context.Context, id int) error {
	return s.Repo.DeleteFood(ctx, id)
}

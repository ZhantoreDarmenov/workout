package services

import (
	"context"

	"workout/internal/models"
	"workout/internal/repositories"
)

// ExerciseService provides business logic for exercises.
type ExerciseService struct {
	Repo *repositories.ExerciseRepository
}

func (s *ExerciseService) CreateExercise(ctx context.Context, ex models.Exercises) (models.Exercises, error) {
	return s.Repo.CreateExercise(ctx, ex)
}

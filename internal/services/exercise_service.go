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

func (s *ExerciseService) UpdateExercise(ctx context.Context, ex models.Exercises) (models.Exercises, error) {
	return s.Repo.UpdateExercise(ctx, ex)
}

func (s *ExerciseService) DeleteExercise(ctx context.Context, id int) error {
	return s.Repo.DeleteExercise(ctx, id)
}

package services

import (
	"context"

	"workout/internal/models"
	"workout/internal/repositories"
)

// ProgramService handles business logic for workout programs.
type ProgramService struct {
	Repo *repositories.ProgramRepository
}

func (s *ProgramService) CreateProgram(ctx context.Context, p models.WorkOutProgram) (models.WorkOutProgram, error) {
	return s.Repo.CreateProgram(ctx, p)
}

func (s *ProgramService) ProgramsByTrainer(ctx context.Context, trainerID int) ([]models.WorkOutProgram, error) {
	return s.Repo.GetProgramsByTrainer(ctx, trainerID)
}

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

func (s *ProgramService) ProgramByID(ctx context.Context, id int) (models.WorkOutProgram, error) {
	return s.Repo.GetProgramByID(ctx, id)
}

func (s *ProgramService) UpdateProgram(ctx context.Context, p models.WorkOutProgram) (models.WorkOutProgram, error) {
	return s.Repo.UpdateProgram(ctx, p)
}

func (s *ProgramService) DeleteProgram(ctx context.Context, id int) error {
	return s.Repo.DeleteProgram(ctx, id)
}

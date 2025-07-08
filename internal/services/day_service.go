package services

import (
	"context"
	"workout/internal/models"
	"workout/internal/repositories"
)

// DayService contains business logic for workout days and progress.
type DayService struct {
	Repo *repositories.DayRepository
}

func (s *DayService) GetDay(ctx context.Context, programID, dayNumber int) (models.DayDetails, error) {
	return s.Repo.GetDayDetails(ctx, programID, dayNumber)
}

func (s *DayService) CompleteDay(ctx context.Context, clientID, dayID int) (models.ProgramProgress, error) {
	return s.Repo.MarkDayCompleted(ctx, clientID, dayID)
}

func (s *DayService) CreateDay(ctx context.Context, day models.Days) (models.Days, error) {
	return s.Repo.CreateDay(ctx, day)
}

func (s *DayService) DaysByProgram(ctx context.Context, programID int) ([]models.DayDetails, error) {
	return s.Repo.DaysByProgram(ctx, programID)
}

func (s *DayService) UpdateDay(ctx context.Context, day models.Days) (models.Days, error) {
	return s.Repo.UpdateDay(ctx, day)
}


func (s *DayService) DeleteDay(ctx context.Context, id int) error {
	return s.Repo.DeleteDay(ctx, id)
}



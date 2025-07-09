package services

import (
	"context"

	"workout/internal/models"
	"workout/internal/repositories"
)

// AnalyticsService wraps business logic for trainer analytics.
type AnalyticsService struct {
	Repo *repositories.AnalyticsRepository
}

func (s *AnalyticsService) TrainerAnalytics(ctx context.Context, trainerID int) (models.TrainerAnalytics, error) {
	return s.Repo.TrainerAnalytics(ctx, trainerID)
}
